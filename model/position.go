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
	EnterAvgPriceGross		float64				// Enter AvgPrice excluding OpenFillFees["totalFees"]
	EnterAvgPriceNet		float64				// Enter AvgPrice including OpenFillFees["totalFees"]
	EnterFillValueNet		float64				// abs(Quantity) * EnterAvgPriceNet

	ExitFillFees 			map[string]float64	// map[feeType]feeAmount
	ExitAvgPriceGross		float64				// Exit AvgPrice excluding ExitFillFees["totalFees"]
	ExitAvgPriceNet			float64				// Exit AvgPrice including ExitFillFees["totalFees"]
	ExitFillValueNet		float64				// abs(Quantity) * ExitAvgPriceNet

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
	p.EnterAvgPriceNet = p.EnterAvgPriceGross + (p.EnterFillFees["TotalFees"] / math.Abs(fill.Quantity)) // When Enter fees make it cost more
	p.EnterFillValueNet = math.Abs(fill.Quantity) * p.EnterAvgPriceNet

	// Exit Fees
	p.ExitFillFees = make(map[string]float64)
	p.ExitFillFees["ExchangeFee"] = 0.0
	p.ExitFillFees["SlippageFee"] = 0.0
	p.ExitFillFees["NetworkFee"] = 0.0
	p.ExitFillFees["TotalFees"] = 0.0

	// Exit Price & Value
	p.ExitAvgPriceGross = 0.0
	p.ExitAvgPriceNet = 0.0
	p.ExitFillValueNet = 0.0

	// Current Symbol Price & Position Value
	p.CurrentSymbolPrice = p.EnterAvgPriceGross  // approx
	p.CurrentMarketValue = p.EnterFillValueNet   // approx

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
	unrealExitFillValue := p.CurrentMarketValue - p.EnterFillFees["TotalFees"]  // Approximate exit fees with enter fees
	unrealProfitLoss, err := calculateProfitLoss(p.Direction, p.Quantity, unrealExitFillValue, p.EnterFillValueNet)
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
	// Todo: Is this fees calculation on ExitAvgPriceNet correct? Does it feed into profitLoss calc correctly for short and long?
	p.ExitAvgPriceGross = fill.FillValueGross / math.Abs(fill.Quantity)
	p.ExitAvgPriceNet = p.ExitAvgPriceGross - (p.ExitFillFees["TotalFees"] / math.Abs(fill.Quantity)) // When Exit fees make it less valuable
	p.ExitFillValueNet = math.Abs(fill.Quantity) * p.ExitAvgPriceNet

	// Result Profit & Loss
	//resultProfitLoss, err := calculateProfitLoss(p.Direction, p.Quantity, p.ExitFillValueNet, p.EnterFillValueNet)
	resultProfitLoss, err := calculateProfitLossV2(*p)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed Position.Exit() for position: %+v", p))
	}
	p.ResultProfitLoss = resultProfitLoss

	return nil
}

// Todo: https://help.bybit.com/hc/en-us/articles/900000630066-P-L-calculations-USDT-Contract-
func calculateProfitLossV2(position Position) (float64, error) {
	var profitLoss float64

	if position.Direction == DirectionLong && position.Quantity > 0 {

		profitLoss = (position.ExitAvgPriceNet * math.Abs(position.Quantity)) + position.EnterFillValueNet

	} else if position.Direction == DirectionShort && position.Quantity < 0 {


		profitLoss = position.EnterFillValueNet - (position.ExitAvgPriceNet * math.Abs(position.Quantity))

	} else {
		return profitLoss, errors.New("failed calculateProfitLoss due to ambiguous Direction & Quantity")
	}

	return profitLoss, nil

}

// calculateProfitLoss calculates the Unreal or Result Profit&Loss given the enter & exit context
func calculateProfitLoss(direction string, quantity float64, exitFillValue float64, enterFillValue float64) (float64, error) {
	var profitLoss float64
	if direction == DirectionLong && quantity > 0 {
		profitLoss = exitFillValue - enterFillValue
	} else if direction == DirectionShort && quantity < 0 {
		profitLoss = enterFillValue - exitFillValue
	} else {
		return profitLoss, errors.New("failed calculateProfitLoss due to ambiguous Direction & Quantity")
	}
	return profitLoss, nil
}