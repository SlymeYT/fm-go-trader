package data

import (
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sheerun/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/config"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
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

// UpdateData updates the currentSymbolData field & adds a MarketEvent to the queue to notify Strategy & Portfolio
func (sh *historicHandler) UpdateData() {
	// Increment latest bar index
	sh.latestBarIndex++

	// Add latest bar to currentSymbolData
	latestBar := model.Bar{
		Timestamp: sh.allSymbolData.Timestamps[sh.latestBarIndex],
		Open:      sh.allSymbolData.Opens[sh.latestBarIndex],
		High:      sh.allSymbolData.Highs[sh.latestBarIndex],
		Low:       sh.allSymbolData.Lows[sh.latestBarIndex],
		Close:     sh.allSymbolData.Closes[sh.latestBarIndex],
		Volume:    sh.allSymbolData.Volumes[sh.latestBarIndex],
	}
	sh.currentSymbolData.AddBar(latestBar)

	// Add MarketEvent to the queue
	sh.eventQ.Append(model.MarketEvent{
		TraceId: uuid.New(),
		Timestamp: latestBar.Timestamp,
		Symbol: sh.symbol,
		Close: latestBar.Close,
	})
}

// GetLatestData returns a tuple of (data up to the current timestamp, latest bar index)
func (sh *historicHandler) GetLatestData() (*model.SymbolData, int64) {
	return &sh.currentSymbolData, sh.latestBarIndex
}

// NewHistoricHandler returns an instance of a data.historicHandler
func NewHistoricHandler(cfg config.Trader, log *zap.Logger,  eventQ *queue.Queue) (*historicHandler, error) {
	filePath := buildCSVFilePath(cfg)
	log.Debug(fmt.Sprintf("loading CSV symbol data with file path: %s", filePath))

	allSymbolData, err := loadCSVSymbolData(filePath)
	if err != nil {
		return &historicHandler{}, errors.Wrap(err, "failed to load CSV data")
	}

	var latestBarIndex int64 = -1
	var currentSymbolData model.SymbolData

	handler := &historicHandler{
		log:            	log,
		eventQ:         	eventQ,
		symbol:         	cfg.Symbol,
		allSymbolData:  	allSymbolData,
		currentSymbolData:	currentSymbolData,
		latestBarIndex: 	latestBarIndex,
	}

	return handler, nil
}

// buildCSVFilePath returns a file path string in the format "dataDirectory + symbol + _ + timeframe + fileExtension"
func buildCSVFilePath(cfg config.Trader) string {
	return fmt.Sprintf("%s%s_%s.csv", dataDirectory, cfg.Symbol, cfg.Timeframe)
}

// loadCSVSymbolData
func loadCSVSymbolData(filePath string) (model.SymbolData, error) {
	lines, err := ReadCSV(filePath)
	if err != nil {
		return model.SymbolData{}, err
	}

	// Loop through (ignoring headers) & build arrays for symbolData struct
	var timestamps []time.Time
	var opens []float64
	var highs []float64
	var lows []float64
	var closes []float64
	var volumes []uint64

	for index, line := range lines[1:] {
		// Add +1 to index to reflect true CSV line number for logging
		index++
		// Timestamp
		timestamp, err := time.Parse(timestampLayoutIso, line[0])
		if err != nil {
			return model.SymbolData{}, errors.Wrap(err, fmt.Sprintf("failed to parse timestamp at index %v", index))
		}
		timestamps = append(timestamps, timestamp)
		// Open
		open, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return model.SymbolData{}, errors.Wrap(err, fmt.Sprintf("failed to parse open at index %v", index))
		}
		opens = append(opens, open)
		// High
		high, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return model.SymbolData{}, errors.Wrap(err, fmt.Sprintf("failed to parse high at index %v", index))
		}
		highs = append(highs, high)
		// Low
		low, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			return model.SymbolData{}, errors.Wrap(err, fmt.Sprintf("failed to parse low at index %v", index))
		}
		lows = append(lows, low)
		// Close
		adjClose, err := strconv.ParseFloat(line[5], 64)
		if err != nil {
			return model.SymbolData{}, errors.Wrap(err, fmt.Sprintf("failed to parse close at index %v", index))
		}
		closes = append(closes, adjClose)
		// Volume
		volume, err := strconv.ParseUint(line[6], 10,64)
		if err != nil {
			return model.SymbolData{}, errors.Wrap(err, fmt.Sprintf("failed to parse volume at index %v", index))
		}
		volumes = append(volumes, volume)
	}

	allSymbolData := model.SymbolData{
		Timestamps: timestamps,
		Opens:      opens,
		Highs:      highs,
		Lows:       lows,
		Closes:     closes,
		Volumes:    volumes,
		Indicators: make(map[string][]interface{}),
	}

	return allSymbolData, nil
}

// ReadCSV reads the file at the provided filePath and returns a 2D array of strings to represent its contents
func ReadCSV(filePath string) ([][]string, error) {
	// Todo: Improve parser w/ https://github.com/dirkolbrich/gobacktest/blob/668a578c68f771714e56e6bedb9744068b6ee2d2/data/data-csv.go#L144
	// Open the file
	csvFile, err := os.Open(filePath)
	if err != nil {
		return [][]string{}, err
	}
	defer csvFile.Close()

	// Read file
	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}