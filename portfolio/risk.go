package portfolio

import "gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"

type RiskManager interface {
	EvaluateOrder(model.OrderEvent) (*model.OrderEvent, error)
}

type Risk struct {
}

// EvaluateOrder manages the risk of an order by refining it, or cancelling it
func (r *Risk) EvaluateOrder(order model.OrderEvent) (*model.OrderEvent, error) {
	return &order, nil
}