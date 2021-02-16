package execution

import (
	"github.com/eapache/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/config"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"go.uber.org/zap"
	"math"
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
		Decision: order.Decision,
	}
	fill.ExchangeFee = fill.CalculateExchangeFee() 		 // 0.0
	fill.SlippageFee = fill.CalculateSlippageFee()		 // 0.0
	fill.NetworkFee = fill.CalculateNetworkFee()		 // 0.0

	// Since simulatedExecution, approximate FillValueGross with most recent market close
	fill.FillValueGross = math.Abs(fill.Quantity) * order.Close //fill.CalculateFillValueGross()

	se.eventQ.Add(fill)
	return nil
}

// NewSimulatedExecution constructs an Execution instance
func NewSimulatedExecution(cfg config.Trader, eventQ *queue.Queue) *simulatedExecution {
	return &simulatedExecution{
		log:      cfg.Log,
		eventQ:   eventQ,
		exchange: cfg.Exchange,
	}
}