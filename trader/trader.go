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
	for {
		// Todo: Need to add check for stop loss and take profit - probably in strategy
		if t.data.ShouldContinue() {
			t.data.UpdateData()
		} else {
			t.log.Info("Backtest has finished.")
			_ = t.DisplayResults()
			// Save & Print results
			// Reset trader instance ready for another run
			break
		}

		for {
			if t.eventQ.Length() > 0 {
				e := t.eventQ.Get(0) // 0 or -1?
				t.eventQ.Remove()
				switch e.(type) {
				case model.MarketEvent:
					repr, _ := json.Marshal(e.(model.MarketEvent))
					t.log.Info(fmt.Sprintf("MARKET: %s", string(repr)))
					err := t.strategy.GenerateSignal(e.(model.MarketEvent))
					if err != nil {
						return errors.Wrap(err, "failed to GenerateSignal()")
					}
					err = t.portfolio.UpdateFromMarket(e.(model.MarketEvent))
					if err != nil {
						return err
					}
				case model.SignalEvent:
					repr, _ := json.Marshal(e.(model.SignalEvent))
					t.log.Info(fmt.Sprintf("SIGNAL: %s", repr))
					err := t.portfolio.GenerateOrders(e.(model.SignalEvent))
					if err != nil {
						return err
					}
				case model.OrderEvent:
					repr, _ := json.Marshal(e.(model.OrderEvent))
					t.log.Info(fmt.Sprintf("ORDER: %s", repr))
					err := t.execution.GenerateFills(e.(model.OrderEvent))
					if err != nil {
						return err
					}
				case model.FillEvent:
					repr, _ := json.Marshal(e.(model.FillEvent))
					t.log.Info(fmt.Sprintf("FILL: %s", repr))
					err := t.portfolio.UpdateFromFill(e.(model.FillEvent))
					if err != nil {
						return err
					}
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