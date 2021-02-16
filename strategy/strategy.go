package strategy

import (
	"github.com/eapache/queue"
	"github.com/markcheno/go-talib"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/config"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/data"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"go.uber.org/zap"
	"time"
)

type Strategy interface {
	GenerateSignal(model.MarketEvent) error
}

type rsiStrategy struct {
	log          *zap.Logger
	eventQ       *queue.Queue
	data         data.Handler
	symbol 		 string
}

// GenerateSignal analyses the current symbol data and appends a
func (s *rsiStrategy) GenerateSignal(market model.MarketEvent) error {
	// Todo: Add some data validation here? Or perhaps return an error from GetLatestBar if data is shit

	// Get current available data and the index of the latest bar
	currentData, latestBarIndex := s.data.GetLatestData()

	// Calculate RSI array
	rsiPeriod := 2
	if latestBarIndex < int64(rsiPeriod) {
		return nil
	}
	rsi2Array := talib.Rsi(currentData.Closes, 2)

	// Construct SignalPairs map
	signalPairs := make(map[string]float32)
	if rsi2Array[latestBarIndex] < 40 {
		signalPairs[model.DecisionLong] = determineSignalStrength()
	}
	if rsi2Array[latestBarIndex] > 60 {
		signalPairs[model.DecisionCloseLong] = determineSignalStrength()
	}
	if rsi2Array[latestBarIndex] > 60 {
		signalPairs[model.DecisionShort] = determineSignalStrength()
	}
	if rsi2Array[latestBarIndex] < 40 {
		signalPairs[model.DecisionCloseShort] = determineSignalStrength()
	}

	// If any SignalPairs
	if len(signalPairs) != 0 {
		// Append SignalEvent to the queue
		s.eventQ.Add(model.SignalEvent{
			EventType: model.EventSignal,
			TraceId: 	 market.TraceId,
			Timestamp:   time.Now().Truncate(time.Nanosecond),
			Symbol:      s.symbol,
			SignalPairs: signalPairs,
		})
	}

	return nil
}

// NewSimpleRSIStrategy constructs a new Strategy instance
func NewSimpleRSIStrategy(cfg config.Trader, eventQ *queue.Queue, data data.Handler) *rsiStrategy {
	return &rsiStrategy{
		log:    cfg.Log,
		eventQ: eventQ,
		data:   data,
		symbol: cfg.Symbol,
	}
}

// determineSignalStrength calculates the strength of a signal advise
func determineSignalStrength() float32{
	return 1.0
}