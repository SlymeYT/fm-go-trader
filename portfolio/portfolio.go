package portfolio

import "gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"

type Portfolio interface {
	UpdateFromMarket() error
	GenerateOrders(model.SignalEvent) error
	UpdateFromFill(model.FillEvent) error
}