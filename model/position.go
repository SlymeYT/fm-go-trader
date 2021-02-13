package model

import (
	"github.com/google/uuid"
	"time"
)

type Position struct {
	LastUpdateTraceId		uuid.UUID
	LastUpdateTimestamp 	time.Time
	Symbol 					string
	Quantity 				float64				// abs(Quantity)

	EnterFillFees			map[string]float64 	// map[feeType]fee
	EnterAvgPriceGross		float64				// Enter AvgPrice excluding OpenFillFees["totalFees"]
	EnterAvgPriceNet		float64				// Enter AvgPrice including OpenFillFees["totalFees"]
	EnterFillValue			float64				// Quantity * EnterAvgPriceNet

	ExitFillFees 			map[string]float64	// map[feeType]fee
	ExitAvgPriceGross		float64				// Exit AvgPrice excluding ExitFillFees["totalFees"]
	ExitAvgPriceNet			float64				// Exit AvgPrice including ExitFillFees["totalFees"]
	ExitFillValue			float64				// Quantity * ExitAvgPriceNet

	CurrentSymbolPrice 		float64				// Symbol current close price
	CurrentMarketValue 		float64				// abs(Quantity) * CurrentSymbolPrice

	UnrealProfitLoss		float64 			// unrealised P&L whilst Position open
	ResultProfitLoss		float64 			// realised P&L after Position closed
}

func (p *Position) Update(market MarketEvent) {
	p.LastUpdateTraceId = market.TraceId
	p.LastUpdateTimestamp = market.Timestamp

	p.CurrentSymbolPrice = market.Close
	p.CurrentMarketValue = market.Close * p.Quantity

	// Unreal Profit & Loss
	// Todo: Work out how to determine unreal PnL, what to use for direction, quantity, and therefore how to calculate profits!
}

func (p *Position) Enter(fill FillEvent) {

}

func (p *Position) Exit(fill FillEvent) {

}

