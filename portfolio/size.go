package portfolio

import (
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"math"
)

type SizeManager interface {
	SizeOrder(*model.OrderEvent, float32, model.Position, float64) error
}

type Size struct {
	DefaultOrderValue 	float64
}

func (s *Size) SizeOrder(order *model.OrderEvent, decisionStrength float32, position model.Position, price float64) error {
	strength := float64(decisionStrength)

	// If order is an exit
	if order.IsExit() {
		enterQuantity := position.Quantity 	// +ve or -ve Quantity depending on Direction
		order.Quantity = (0.0 - enterQuantity) * strength
	}

	// If order is an entry
	defaultOrderSize := math.Floor(s.DefaultOrderValue / price)
	if order.IsLong() {
		order.Quantity = defaultOrderSize * strength
	}
	if order.IsShort() {
		order.Quantity = -defaultOrderSize * strength
	}

	return nil
}