package trader

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sheerun/queue"
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
		// Todo: Need to add check for stop loss and take profit
		if t.data.ShouldContinue() {
			t.data.UpdateData()
		} else {
			t.log.Info("Backtest has finished.")
			// Save & Print results
			// Reset trader instance ready for another run
			break
		}

		for {
			if t.eventQ.Length() > 0 {
				e := t.eventQ.Pop()
				switch e.(type) {
				case model.MarketEvent:
					t.log.Info(fmt.Sprintf("%+v", e.(model.MarketEvent)))
					_ = t.strategy.GenerateSignal(e.(model.MarketEvent))
					_ = t.portfolio.UpdateFromMarket(e.(model.MarketEvent))
				case model.SignalEvent:
					t.log.Info(fmt.Sprintf("%+v", e.(model.SignalEvent)))
					_ = t.portfolio.GenerateOrders(e.(model.SignalEvent))
				case model.OrderEvent:
					t.log.Info(fmt.Sprintf("%+v", e.(model.OrderEvent)))
					_ = t.execution.GenerateFills(e.(model.OrderEvent))
				case model.FillEvent:
					t.log.Info(fmt.Sprintf("%+v", e.(model.FillEvent)))
					_ = t.portfolio.UpdateFromFill(e.(model.FillEvent))
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