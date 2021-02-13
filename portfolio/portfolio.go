package portfolio

import (
	"github.com/sheerun/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/data"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"go.uber.org/zap"
)

type Portfolio interface {
	UpdateFromMarket() error
	GenerateOrders(model.SignalEvent) error
	UpdateFromFill(model.FillEvent) error
}

type portfolio struct {
	log              *zap.Logger
	eventQ           *queue.Queue
	data             data.Handler
	symbol           string
	initialCash      float64
	currentCash      float64
	currentValue 	 float64
	orders           []model.OrderEvent
	fills		     []model.FillEvent
	holdings         map[string]model.Position
	historicHoldings map[string][]model.Position
}