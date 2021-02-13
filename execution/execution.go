package execution

import (
	"github.com/sheerun/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/config"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"go.uber.org/zap"
	"time"
)

type Execution interface {
	GenerateFills(model.OrderEvent) error
}

type simulatedExecution struct {
	log      *zap.Logger
	eventQ   *queue.Queue
	exchange string
}

// GenerateFills takes an OrderEvent, executes it, and produces a FillEvent that is appended to the event queue
func (se *simulatedExecution) GenerateFills(order model.OrderEvent) error {
	// Todo: Add latency, slippage, etc

	// Assume all orders are filled at the market price
	fill := model.FillEvent{
		TraceId: order.TraceId,
		Timestamp: time.Now().Truncate(time.Nanosecond),
		Symbol:    order.Symbol,
		Exchange:  se.exchange,
		Quantity:  order.Quantity,
		Direction: order.Direction,
	}
	fill.ExchangeFee = fill.CalculateExchangeFee() 		// 0.0
	fill.SlippageFee = fill.CalculateSlippageFee()		// 0.0
	fill.NetworkFee = fill.CalculateNetworkFee()		// 0.0
	fill.FillCost = fill.CalculateFillCost() 			// 0.0

	se.eventQ.Append(fill)
	return nil
}

// NewSimulatedExecution constructs an Execution instance
func NewSimulatedExecution(cfg config.Trader, logger *zap.Logger, eventQ *queue.Queue) *simulatedExecution {
	return &simulatedExecution{
		log:      logger,
		eventQ:   eventQ,
		exchange: cfg.Exchange,
	}
}
