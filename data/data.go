package data

import (
	"github.com/sheerun/queue"
	"go.uber.org/zap"
	"time"
)

const(
	dataDirectory = "data/"
	timestampLayoutIso = "2006-01-02"
)

type Handler interface {
	ShouldContinue() bool
	UpdateData() error
	GetLatestData() (*symbolData, int64)
}

type Config struct {
	Log             *zap.Logger
	EventQ          *queue.Queue
	Ticker          string
	Timeframe       string
	Exchange        string
}

// historicHandler is a Handler for backtesting trading strategies with historic data
type historicHandler struct {
	log               *zap.Logger	// Pointer to universal logger
	eventQ            *queue.Queue	// Pointer to the trader pair's event queue
	symbol            string     	// symbol the data is representing
	allSymbolData     symbolData 	// All the data available from historic data file
	currentSymbolData symbolData 	// Data available up to current timestamp
	latestBarIndex    int64      	// Current index of the latest bar in the symbolData
}

// symbolData represents a symbol's struct of market data arrays (OHLCV) and associated indicators values
type symbolData struct {
	Timestamps 	[]time.Time
	Opens 		[]float64
	Highs 		[]float64
	Lows 		[]float64
	Closes 		[]float64
	Volumes 	[]uint64
	Indicators 	map[string][]interface{}
}

// ShouldContinue determines if the market data feed should be terminated
func (sh *historicHandler) ShouldContinue() bool {
	var shouldContinue bool
	if sh.latestBarIndex < int64(len(sh.allSymbolData.Timestamps)) - 1 {
		shouldContinue = true
	}
	return shouldContinue
}



