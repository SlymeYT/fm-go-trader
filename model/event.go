package model

import (
	"github.com/google/uuid"
	"time"
)

type MarketEvent struct {
	TraceId 	uuid.UUID
	Timestamp 	time.Time
}

type SignalEvent struct {
	TraceId 	uuid.UUID
	Timestamp 	time.Time
	Symbol 		string
	SignalPairs map[string]float32 	// map[direction]strength
}

type OrderEvent struct {
	TraceId 	uuid.UUID
	Timestamp 	time.Time
	Symbol    	string
	Quantity   	float64
	Direction  	string
}

type FillEvent struct {
	TraceId 		uuid.UUID
	Timestamp  		time.Time
	Symbol     		string
	Exchange   		string
	Quantity   		float64
	Direction  		string
	FillCost   		float64
	CommissionFee 	float64
	SlippageFee		float64
}

func (f *FillEvent) CalculateCommissionFee() float64 {
	return 0.0
}

func (f *FillEvent) CalculateSlippageFee() float64 {
	return 0.0
}

func (f *FillEvent) CalculateFillCost() float64 {
	return 0.0
}