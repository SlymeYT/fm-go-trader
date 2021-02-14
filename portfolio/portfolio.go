package portfolio

import (
	"github.com/pkg/errors"
	"github.com/sheerun/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/config"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/data"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/strategy"
	"go.uber.org/zap"
)

type Portfolio interface {
	UpdateFromMarket(event model.MarketEvent) error
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

func (p *portfolio) UpdateFromMarket(market model.MarketEvent) error {
	// Update currentHoldings
	if position, isInvested := p.isInvested(p.symbol); isInvested {
		err := position.Update(market)
		if err != nil {
			return errors.Wrap(err, "failed portfolio.UpdateFromMarket()")
		}
		p.holdings[p.symbol] = position
	}

	return nil
}

func (p *portfolio) GenerateOrders(signal model.SignalEvent) error {
	// Todo:
	//  - Check if we are invested
	//  - If not, check if we have any available cash
	//  - Parse signal.signalPairs map to find decision to action & associated strength
	//  - Pass to risk manager to refine or cancel order
	//  - Pass to size manager to set order quantity / value
	//  - Add order to orders book
	//  - Append order to event queue

	// SignalEvent is for a Symbol we are already invested in
	_, isInvested := p.isInvested(signal.Symbol)



	return nil
}

func (p *portfolio) UpdateFromFill(fill model.FillEvent) error {
	return nil
}

func (p *portfolio) isInvested(symbol string) (model.Position, bool) {
	position, isInHoldings := p.holdings[symbol]
	// If present in current holdings & exit fill value is zero
	if isInHoldings && position.ExitFillValueNet == 0 {
		return position, true
	}
	return position, false
}

func (p *portfolio) parseSignalDecisions(signalPairs map[string]float32) (strength float32, decision string) {
	//strengthLong, adviseLong := signalPairs[strategy.DecisionLong]
	//strengthCloseLong, adviseCloseLong := signalPairs[strategy.DecisionCloseLong]
	//strengthShort, adviseShort := signalPairs[strategy.DecisionShort]
	//strengthCloseShort, adviseCloseShort := signalPairs[strategy.DecisionCloseShort]

	return
}

func NewPortfolio(cfg config.Trader, eventQ *queue.Queue, data data.Handler) *portfolio {
	return &portfolio{
		log:              cfg.Log,
		eventQ:           eventQ,
		data:             data,
		sizeManager:      &Size{DefaultOrderValue: cfg.DefaultOrderValue},
		riskManager:      &Risk{DefaultOrderType: OrderTypeMarket},
		symbol:           cfg.Symbol,
		initialCash:      cfg.StartingCash,
		currentCash:      cfg.StartingCash,
		currentValue:     cfg.StartingCash,
		orders:           []model.OrderEvent{},
		fills:            []model.FillEvent{},
		holdings:         make(map[string]model.Position),
		historicHoldings: make(map[string][]model.Position),
	}
}