package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

// MarketEvent (data) is the system heartbeat & represents the arrival of new data for the strategy to interpret
type MarketEvent struct {
	TraceId 	uuid.UUID
	Timestamp 	time.Time
	Symbol		string
	Close		float64
}

// SignalEvent (strategy) are advisory signals for the portfolio to interpret
type SignalEvent struct {
	TraceId 	uuid.UUID
	Timestamp 	time.Time
	Symbol 		string
	SignalPairs map[string]float32 	// map[Decision]Strength
}

// OrderEvent (portfolio) are actions for the execution handler to execute
type OrderEvent struct {
	TraceId 	uuid.UUID
	Timestamp 	time.Time
	Symbol    	string
	OrderType 	string  	// MARKET, LIMIT etc
	Quantity   	float64		// +ve or -ve Quantity depending on Decision
	Decision  	string		// LONG, CLOSE_LONG, SHORT or CLOSE_SHORT
}

// FillEvent (execution) are journals of work done sent back to the portfolio to interpret and update holdings
type FillEvent struct {
	TraceId 		uuid.UUID
	Timestamp  		time.Time
	Symbol     		string
	Exchange   		string
	Quantity   		float64		// +ve or -ve Quantity depending on Decision
	Decision  		string 		// LONG, CLOSE_LONG, SHORT or CLOSE_SHORT
	FillValueGross  float64		// Todo: say that fees are not included
	ExchangeFee 	float64
	SlippageFee		float64
	NetworkFee		float64
}

// DetermineFillDirection determines the Direction of a FillEvent based on it's Quantity and Decision
func (f *FillEvent) DetermineFillDirection() (string, error) {
	var direction string
	if (f.Decision == "LONG" || f.Decision == "CLOSE_LONG") && f.Quantity > 0 {
		direction = "LONG"
	} else if (f.Decision == "SHORT" || f.Decision == "CLOSE_SHORT") && f.Quantity < 0{
		direction = "SHORT"
	} else {
		return direction, errors.New(fmt.Sprintf("failed FillEvent.DetermineFillDirection() due to ambiguous Quanity & Decision, FillEvent: %+v", f))
	}
	return direction, nil
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