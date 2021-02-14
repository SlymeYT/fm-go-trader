package portfolio

import "gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"

const (
	OrderTypeMarket = "MARKET"
)

type RiskManager interface {
	EvaluateOrder(*model.OrderEvent) error
}

type Risk struct {
	DefaultOrderType string
}

// EvaluateOrder manages the risk of an order by refining it, or cancelling it
func (r *Risk) EvaluateOrder(order *model.OrderEvent) error {
	order.OrderType = OrderTypeMarket
	return nil
}