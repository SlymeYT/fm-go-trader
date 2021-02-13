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
	sizeManager		 SizeManager
	riskManager		 RiskManager
	symbol           string
	initialCash      float64
	currentCash      float64
	currentValue 	 float64
	orders           []model.OrderEvent
	fills		     []model.FillEvent
	holdings         map[string]model.Position
	historicHoldings map[string][]model.Position
}

func (p *portfolio) UpdateFromMarket() error {
	return nil
}

func (p *portfolio) GenerateOrders() error {
	return nil
}

func (p *portfolio) UpdateFromFill() error {
	return nil
}

func NewPortfolio() *portfolio {
	return &portfolio{}
}

func parseAdvise(signalPairs map[string]float32) string {
	return "LONG" // SHORT or EXIT
}