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