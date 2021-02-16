package trader

import (
	"encoding/json"
	"fmt"
	"github.com/eapache/queue"
	"github.com/pkg/errors"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/config"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/data"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/execution"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/portfolio"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/strategy"
	"go.uber.org/zap"
	"os"
	"time"
)

const (
	resultDirectory = "data/result/"
)

type Trader interface {
	Run() error
	DisplayResults() error
}

type trader struct {
	log       *zap.Logger
	eventQ    *queue.Queue
	data      data.Handler
	strategy  strategy.Strategy
	portfolio portfolio.Portfolio
	execution execution.Execution
}

func (t *trader) Run() error {
	// Setup Market Event Stream File
	encoder, encoderClean, err := setupEventLog("ETH-USD")
	if err != nil {
		return err
	}

	for {
		// Todo: Need to add check for stop loss and take profit - probably in strategy
		if t.data.ShouldContinue() {
			t.data.UpdateData()
		} else {
			t.log.Info("Backtest has finished.")
			_ = t.DisplayResults()
			err := encoderClean()
			if err != nil {
				return err
			}
			// Save & Print results
			// Reset trader instance ready for another run
			break
		}

		for {
			if t.eventQ.Length() > 0 {
				e := t.eventQ.Get(0)
				t.eventQ.Remove()

				switch e.(type) {
				case model.MarketEvent:
					err := t.strategy.GenerateSignal(e.(model.MarketEvent))
					if err != nil {
						return errors.Wrap(err, "failed to GenerateSignal()")
					}
					err = t.portfolio.UpdateFromMarket(e.(model.MarketEvent))
					if err != nil {
						return errors.Wrap(err, "failed to UpdateFromMarket()")
					}
					_ = encoder.Encode(e.(model.MarketEvent))
				case model.SignalEvent:
					err := t.portfolio.GenerateOrders(e.(model.SignalEvent))
					if err != nil {
						return errors.Wrap(err, "failed to GenerateOrders()")
					}
					_ = encoder.Encode(e.(model.SignalEvent))
				case model.OrderEvent:
					err := t.execution.GenerateFills(e.(model.OrderEvent))
					if err != nil {
						return errors.Wrap(err, "failed to GenerateFills()")
					}
					_ = encoder.Encode(e.(model.OrderEvent))
				case model.FillEvent:
					err := t.portfolio.UpdateFromFill(e.(model.FillEvent))
					if err != nil {
						return errors.Wrap(err, "failed to UpdateFromFill()")
					}
					_ = encoder.Encode(e.(model.FillEvent))
				}
			} else {
				// Inner loop would break when the event queue is empty and we need another data drop
				break
			}
		}
		// This is the heartbeat -> would be frequency of poll to get data from execution
		//time.Sleep(2*time.Millisecond)
	}
	return nil
}

func setupEventLog(symbol string) (*json.Encoder, func() (err error), error) {
	// Create *File (io.Writer)
	fileName := fmt.Sprintf("%seventLog_%s_%v", resultDirectory, symbol, time.Now().Truncate(time.Millisecond))
	file, err := os.Create(fileName)
	if err != nil {
		return nil, nil, err
	}

	// Create Json Encoder w/ *File
	encoder := json.NewEncoder(file)

	// Create cleanup func
	eventLogClose := func() (err error) {
		return file.Close()
	}

	return encoder, eventLogClose, nil
}

func (t *trader) DisplayResults() error {
	initialCash, currentCash, currentValue, positions := t.portfolio.GetPortfolio()

	fmt.Printf("\nDisplay Results:\n")
	fmt.Printf("Starting Cash: %v\n", initialCash)
	fmt.Printf("Ending Cash: %v\n", currentCash)
	fmt.Printf("Ending Value: %v\n", currentValue)
	fmt.Printf("Number Trades: %v\n", len(positions["ETH-USD"]))

	totalProfit := currentValue - initialCash
	fmt.Printf("Total Profit: %v\n", totalProfit)
	percentProfit := (totalProfit / initialCash) * 100
	fmt.Printf("Total Percent Profit: %v\n", percentProfit)

	return nil
}

func NewTrader(cfg config.Trader) (*trader, error) {

	eventQ := queue.New()

	dataHandler, err := data.NewHistoricHandler(cfg, eventQ)
	if err != nil {
		return &trader{}, errors.Wrap(err, "failed to init dataHandler")
	}
	basicStrategy := strategy.NewSimpleRSIStrategy(cfg, eventQ, dataHandler)
	basicPortfolio := portfolio.NewPortfolio(cfg, eventQ, dataHandler)
	basicExecution := execution.NewSimulatedExecution(cfg, eventQ)

	trader := &trader{
		log:       cfg.Log,
		eventQ:    eventQ,
		data:      dataHandler,
		strategy:  basicStrategy,
		portfolio: basicPortfolio,
		execution: basicExecution,
	}
	return trader, nil
}