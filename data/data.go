package data

import (
	"github.com/sheerun/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"go.uber.org/zap"
)

const(
	dataDirectory = "data/"
	timestampLayoutIso = "2006-01-02"
)

type Handler interface {
	ShouldContinue() bool
	UpdateData() error
	GetLatestData() (*model.SymbolData, int64)
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
	log               *zap.Logger		// Pointer to repository logger
	eventQ            *queue.Queue		// Pointer to the trader pair's event queue
	symbol            string     		// symbol the data is representing
	allSymbolData     model.SymbolData 	// All the data available from historic data file
	currentSymbolData model.SymbolData 	// Data available up to current timestamp
	latestBarIndex    int64      		// Current index of the latest bar in the symbolData
}

// ShouldContinue determines if the market data feed should be terminated
func (sh *historicHandler) ShouldContinue() bool {
	var shouldContinue bool
	if sh.latestBarIndex < int64(len(sh.allSymbolData.Timestamps)) - 1 {
		shouldContinue = true
	}
	return shouldContinue
}

// UpdateData updates the latestTickerData field & adds a MarketEvent to the queue to notify Strategy & Portfolio
func (sh *historicHandler) UpdateData() error {
	// Increment latest bar index
	sh.latestBarIndex++

	// Add latest bar to currentSymbolData
	sh.currentSymbolData.AddBar(model.Bar{
		Timestamp: sh.allSymbolData.Timestamps[sh.latestBarIndex],
		Open: sh.allSymbolData.Opens[sh.latestBarIndex],
		High: sh.allSymbolData.Highs[sh.latestBarIndex],
		Low: sh.allSymbolData.Lows[sh.latestBarIndex],
		Close: sh.allSymbolData.Closes[sh.latestBarIndex],
		Volume: sh.allSymbolData.Volumes[sh.latestBarIndex],
	})

	// Add MarketEvent to the queue
	sh.eventQ.Append(model.MarketEvent{})
	return nil
}

// GetLatestData returns a tuple of (data up to the current timestamp, latest bar index)
func (sh *historicHandler) GetLatestData() (*model.SymbolData, int64) {
	return &sh.currentSymbolData, sh.latestBarIndex
}


