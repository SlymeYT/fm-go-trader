package portfolio

import "gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"

type SizeManager interface {
	SizeOrder(model.OrderEvent) (*model.OrderEvent, error)
}

type Size struct {
	DefaultSize 	float64
	DefaultValue 	float64
}

func (s *Size) SizeOrder(order model.OrderEvent) (*model.OrderEvent, error) {
	order.Quantity = 10
	return &order, nil
}