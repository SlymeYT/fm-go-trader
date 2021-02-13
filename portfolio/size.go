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
	// Todo: Implement a basic version of this so I'm not creating orders that I don't have bean for etc.
	order.Quantity = 10
	return &order, nil
}