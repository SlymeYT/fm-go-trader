package model

import "time"

// Bar represents a symbol's market data state at a fixed interval of time
type Bar struct {
	Timestamp time.Time
	Open float64
	High float64
	Low float64
	Close float64
	Volume uint64
}

// SymbolData represents a symbol's struct of market data arrays (OHLCV) and associated indicators values
type SymbolData struct {
	Timestamps 	[]time.Time
	Opens 		[]float64
	Highs 		[]float64
	Lows 		[]float64
	Closes 		[]float64
	Volumes 	[]uint64
	Indicators 	map[string][]interface{}
}

// AddBar appends each bar field to the relevant SymbolData array
func (td *SymbolData) AddBar(bar Bar) {
	td.Timestamps = append(td.Timestamps, bar.Timestamp)
	td.Opens = append(td.Opens, bar.Open)
	td.Highs = append(td.Highs, bar.High)
	td.Lows = append(td.Lows, bar.Low)
	td.Closes = append(td.Closes, bar.Close)
	td.Volumes = append(td.Volumes, bar.Volume)
}