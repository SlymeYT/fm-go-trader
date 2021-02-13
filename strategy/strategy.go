package strategy

import (
	"github.com/markcheno/go-talib"
	"github.com/sheerun/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/data"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"go.uber.org/zap"
	"time"
)

const (
	DirectionLong 		= "LONG"
	DirectionCloseLong  = "CLOSE_LONG"
	DirectionShort 		= "SHORT"
	DirectionCloseShort = "CLOSE_SHORT"
)

type Strategy interface {
	GenerateSignal() error
}

type simpleRSIStrategy struct {
	log          *zap.Logger
	eventQ       *queue.Queue
	data         data.Handler
	symbol 		 string
}

func (s *simpleRSIStrategy) GenerateSignal() error {
	// Todo: Add some data validation here? Or perhaps return an error from GetLatestBar if data is shit

	// Get current available data and the index of the latest bar
	currentData, latestBarIndex := s.data.GetLatestData()

	// Calculate RSI array
	rsiPeriod := 2
	if latestBarIndex < int64(rsiPeriod) {
		return nil
	}
	rsi2Array := talib.Rsi(currentData.Closes, 2)

	// Construct base SignalEvent
	signalEvent := model.SignalEvent{
		Timestamp: time.Now(),
		Symbol: s.symbol,
	}

	// Construct SignalPairs map
	signalPairs := make(map[string]float32)
	if rsi2Array[latestBarIndex] < 40 {
		signalPairs[DirectionLong] = determineSignalStrength()
	}
	if rsi2Array[latestBarIndex] > 60 {
		signalPairs[DirectionCloseLong] = determineSignalStrength()
	}
	if rsi2Array[latestBarIndex] > 60 {
		signalPairs[DirectionShort] = determineSignalStrength()
	}
	if rsi2Array[latestBarIndex] < 40 {
		signalPairs[DirectionCloseShort] = determineSignalStrength()
	}
	// Add SignalPairs to SignalEvent
	signalEvent.SignalPairs = signalPairs

	s.eventQ.Append(signalEvent)

	return nil
}

func determineSignalStrength() float32{
	return 1.0
}

