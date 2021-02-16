package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math"
	"time"
)

const (
	DirectionLong = "LONG"
	DirectionShort = "SHORT"
)

type Position struct {
	LastUpdateTraceId		uuid.UUID
	LastUpdateTimestamp 	time.Time
	Symbol 					string
	Direction				string				// LONG or SHORT
	Quantity 				float64				// +ve or -ve Quantity of Symbol contracts opened

	EnterFillFees			map[string]float64 	// map[feeType]feeAmount
	EnterAvgPriceGross		float64				// Enter AvgPrice excluding EnterFillFees["totalFees"]
	EnterFillValueGross		float64				// abs(Quantity) * EnterAvgPriceGross

	ExitFillFees 			map[string]float64	// map[feeType]feeAmount
	ExitAvgPriceGross		float64				// Exit AvgPrice excluding ExitFillFees["totalFees"]
	ExitFillValueGross		float64				// abs(Quantity) * ExitAvgPriceGross

	CurrentSymbolPrice 		float64				// Symbol current close price
	CurrentMarketValue 		float64				// abs(Quantity) * CurrentSymbolPrice

	UnrealProfitLoss		float64 			// unrealised P&L whilst Position open
	ResultProfitLoss		float64 			// realised P&L after Position closed
}

// Enter enriches a new Position using information from an enter FillEvent
func (p *Position) Enter(fill FillEvent) error {
	p.LastUpdateTraceId = fill.TraceId
	p.LastUpdateTimestamp = fill.Timestamp
	p.Symbol = fill.Symbol

	// Direction
	direction, err := fill.DetermineFillDirection()
	if err != nil {
		return errors.New("failed to enter new Position")
	}
	p.Direction = direction

	// +ve or -ve Quantity depending on FillEvent Direction
	p.Quantity = fill.Quantity

	// Enter Fees
	p.EnterFillFees = make(map[string]float64)
	p.EnterFillFees["ExchangeFee"] = fill.ExchangeFee
	p.EnterFillFees["SlippageFee"] = fill.SlippageFee
	p.EnterFillFees["NetworkFee"] = fill.NetworkFee
	p.EnterFillFees["TotalFees"] = fill.ExchangeFee + fill.SlippageFee + fill.NetworkFee

	// Enter Price & Value
	p.EnterAvgPriceGross = fill.FillValueGross / math.Abs(fill.Quantity)
	p.EnterFillValueGross = math.Abs(fill.Quantity) * p.EnterAvgPriceGross

	// Exit Fees
	p.ExitFillFees = make(map[string]float64)
	p.ExitFillFees["ExchangeFee"] = 0.0
	p.ExitFillFees["SlippageFee"] = 0.0
	p.ExitFillFees["NetworkFee"] = 0.0
	p.ExitFillFees["TotalFees"] = 0.0

	// Exit Price & Value
	p.ExitAvgPriceGross = 0.0
	//p.ExitAvgPriceNet = 0.0
	p.ExitFillValueGross = 0.0

	// Current Symbol Price & Position Value
	p.CurrentSymbolPrice = p.EnterAvgPriceGross  // approx
	p.CurrentMarketValue = p.EnterFillValueGross   // approx

	// Profit & Loss
	p.UnrealProfitLoss = 0.0
	p.ResultProfitLoss = 0.0

	return nil
}

// Update updates the Position instance on every MarketEvent
func (p *Position) Update(market MarketEvent) error {
	p.LastUpdateTraceId = market.TraceId
	p.LastUpdateTimestamp = market.Timestamp

	p.CurrentSymbolPrice = market.Close
	p.CurrentMarketValue = market.Close * math.Abs(p.Quantity)

	// Unreal Profit & Loss
	unrealProfitLoss, err := calculateUnrealProfitLoss(*p)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed Position.Update() for Position: %+v", p))
	}
	p.UnrealProfitLoss = unrealProfitLoss

	return nil
}

// Exit closes an existing Position with information from an exit FillEvent
func (p *Position) Exit(fill FillEvent) error {
	p.LastUpdateTraceId = fill.TraceId
	p.LastUpdateTimestamp = fill.Timestamp

	// Exit Fees
	p.ExitFillFees["ExchangeFee"] = fill.ExchangeFee
	p.ExitFillFees["SlippageFee"] = fill.SlippageFee
	p.ExitFillFees["NetworkFee"] = fill.NetworkFee
	p.ExitFillFees["TotalFees"] = fill.ExchangeFee + fill.SlippageFee + fill.NetworkFee

	// Exit Price & Value
	p.ExitAvgPriceGross = fill.FillValueGross / math.Abs(fill.Quantity)
	p.ExitFillValueGross = math.Abs(fill.Quantity) * p.ExitAvgPriceGross

	// Result Profit & Loss
	resultProfitLoss, err := calculateResultProfitLoss(*p)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed Position.Exit() for position: %+v", p))
	}
	p.ResultProfitLoss = resultProfitLoss

	return nil
}

// Todo: https://help.bybit.com/hc/en-us/articles/900000630066-P-L-calculations-USDT-Contract-
// calculateUnrealProfitLoss calculates the Unreal Profit&Loss given a copy of the entered Position
func calculateUnrealProfitLoss(position Position) (float64, error) {
	var profitLoss float64

	totalFees := position.EnterFillFees["TotalFees"] * 2
	if position.Direction == DirectionLong && position.Quantity > 0 {
		profitLoss = (position.CurrentMarketValue - position.EnterFillValueGross) - totalFees
	} else if position.Direction == DirectionShort && position.Quantity < 0 {
		profitLoss = (position.EnterFillValueGross - position.CurrentMarketValue) - totalFees
	} else {
		return profitLoss, errors.New("failed calculateUnrealProfitLoss due to ambiguous Direction & Quantity")
	}

	return profitLoss, nil
}

// calculateResultProfitLoss calculates the Result Profit&Loss given a copy of the exited Position
func calculateResultProfitLoss(position Position) (float64, error) {
	var profitLoss float64

	totalFees := position.EnterFillFees["TotalFees"] + position.ExitFillFees["TotalFees"]
	if position.Direction == DirectionLong && position.Quantity > 0 {
		profitLoss = (position.ExitFillValueGross - position.EnterFillValueGross) - totalFees
	} else if position.Direction == DirectionShort && position.Quantity < 0 {
		profitLoss = (position.EnterFillValueGross - position.ExitFillValueGross) - totalFees
	} else {
		return profitLoss, errors.New("failed calculateResultProfitLoss due to ambiguous Direction & Quantity")
	}
	return profitLoss, nil
}