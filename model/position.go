package model

import "time"

type Position struct {
	LastUpdateTimestamp 	time.Time
	Symbol 					string
	Quantity 				int64				// abs(Quantity)

	EntryFillFees			map[string]float64 	// map[feeType]fee
	EntryAvgPriceGross		float64				// Entry AvgPrice excluding OpenFillFees["totalFees"]
	EntryAvgPriceNet		float64				// Entry AvgPrice including OpenFillFees["totalFees"]
	EntryFillValue			float64				// Quantity * EntryAvgPriceNet

	ExitFillFees 			map[string]float64	// map[feeType]fee
	ExitAvgPriceGross		float64				// Exit AvgPrice excluding ExitFillFees["totalFees"]
	ExitAvgPriceNet			float64				// Exit AvgPrice including ExitFillFees["totalFees"]
	ExitFillValue			float64				// Quantity * ExitAvgPriceNet

	CurrentSymbolPrice 		float64		// Symbol current close price
	CurrentMarketValue 		float64		// abs(Quantity) * CurrentSymbolPrice

	UnrealProfitLoss		float64 	// unrealised P&L whilst Position open
	ResultProfitLoss		float64 	// realised P&L after Position closed
}