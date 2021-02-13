package model

import "time"

type MarketEvent struct {
	Timestamp time.Time
}

type SignalEvent struct {
	Timestamp 	time.Time
	Symbol 		string
	SignalPairs map[string]float32 	// map[direction]strength
}

type OrderEvent struct {
	Timestamp 	time.Time
	Symbol    	string
	Quantity   		float64
	Direction  		string
}

type FillEvent struct {
	Timestamp  		time.Time
	Symbol     		string
	Exchange   		string
	Quantity   		float64
	Direction  		string
	FillCost   		float64
	CommissionFee 	float64
	ExchangeFee		float64
	//IsExit 			bool
}

func (f *FillEvent) CalculateCommissionFee() float64 {
	return 0.0
}

func (f *FillEvent) CalculateExchangeFee() float64 {
	return 0.0
}

func (f *FillEvent) CalculateFillCost() float64 {
	return 0.0
}