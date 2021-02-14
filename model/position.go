package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math"
	"time"
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

func (p *Position) Update(market MarketEvent) error {
	p.LastUpdateTraceId = market.TraceId
	p.LastUpdateTimestamp = market.Timestamp

	p.CurrentSymbolPrice = market.Close
	p.CurrentMarketValue = market.Close * math.Abs(p.Quantity)

	// Unreal Profit & Loss
	unrealExitFillValue := p.CurrentMarketValue - p.EnterFillFees["TotalFees"]  // Approximate exit fees with enter fees
	if p.Direction == "LONG" && p.Quantity > 0 {
		p.UnrealProfitLoss = calculateLongProfitLoss(unrealExitFillValue, p.EnterFillValueNet)
	} else if p.Direction == "SHORT" && p.Quantity < 0 {
		p.UnrealProfitLoss = calculateShortProfitLoss(unrealExitFillValue, p.ExitFillValueNet)
	} else {
		return errors.New(fmt.Sprintf("failed Position.Update() due to ambiguous Direction & Quantity. Position: %+v", p))
	}

	return nil
}

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
	p.EnterAvgPriceNet = p.EnterAvgPriceGross + (p.EnterFillFees["TotalFees"] / math.Abs(fill.Quantity))
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

func (p *Position) Exit(fill FillEvent) {

}

func calculateLongProfitLoss(exitFillValue float64, enterFillValue float64) float64 {
	return exitFillValue - enterFillValue
}

func calculateShortProfitLoss(exitFillValue float64, enterFillValue float64) float64 {
	return enterFillValue - exitFillValue
}

