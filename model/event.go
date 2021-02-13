package model

import (
	"github.com/google/uuid"
	"time"
)

// MarketEvent (data) is the system heartbeat & represents the arrival of new data for the strategy to interpret
type MarketEvent struct {
	TraceId 	uuid.UUID
	Timestamp 	time.Time
}

// SignalEvent (strategy) are advisory signals for the portfolio to interpret
type SignalEvent struct {
	TraceId 	uuid.UUID
	Timestamp 	time.Time
	Symbol 		string
	SignalPairs map[string]float32 	// map[direction]strength
}

// OrderEvent (portfolio) are actions for the execution handler to execute
type OrderEvent struct {
	TraceId 	uuid.UUID
	Timestamp 	time.Time
	Symbol    	string
	OrderType 	string  	// MARKET, LIMIT etc
	Quantity   	float64
	Direction  	string
}

// FillEvent (execution) are journals of work done sent back to the portfolio to interpret and update holdings
type FillEvent struct {
	TraceId 		uuid.UUID
	Timestamp  		time.Time
	Symbol     		string
	Exchange   		string
	Quantity   		float64
	Direction  		string
	FillCost   		float64
	ExchangeFee 	float64
	SlippageFee		float64
	NetworkFee		float64
}

// CalculateExchangeFee calculates the exchange fees incurred by the FillEvent
func (f *FillEvent) CalculateExchangeFee() float64 {
	return 0.0
}

// CalculateSlippageFee calculates the slippage fees (losses) incurred by the FillEvent
func (f *FillEvent) CalculateSlippageFee() float64 {
	return 0.0
}

// CalculateNetworkFee calculates the network fees (DEX) incurred by the FillEvent
func (f *FillEvent) CalculateNetworkFee() float64 {
	return 0.0
}

// CalculateFillValue calculates the total value transacted by the FillEvent
func (f *FillEvent) CalculateFillValue() float64 {
	return 0.0
}