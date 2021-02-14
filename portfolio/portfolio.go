package portfolio

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sheerun/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/config"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/data"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"go.uber.org/zap"
	"time"
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
	// Check if the SignalEvent is for a Symbol already invested in
	position, isInvested := p.isInvested(signal.Symbol)

	// If no cash, cannot open a new position -> exit without generating an order
	if !isInvested && p.currentCash == 0.0 {
		return nil
	}

	// Parse SignalPairs map to determine the net OrderEvent decision
	strength, decision := p.parseSignalDecisions(signal.SignalPairs)

	// Construct base OrderEvent
	order := model.OrderEvent{
		TraceId:   signal.TraceId,
		Timestamp: time.Now().Truncate(time.Nanosecond),
		Symbol:    signal.Symbol,
		Decision:  decision,
	}

	// Size order
	// Get current available data and the index of the latest bar
	currentData, latestBarIndex := p.data.GetLatestData()

	err := p.sizeManager.SizeOrder(&order, strength, position, currentData.Closes[latestBarIndex])
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to size order: %+v", order))
	}

	// Manage risk - refine or cancel order
	err = p.riskManager.EvaluateOrder(&order)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to risk evaluate order: %+v", order))
	}

	// Append order to the orders book
	p.orders = append(p.orders, order)

	// Append order to the event queue
	p.eventQ.Append(order)

	return nil
}

func (p *portfolio) UpdateFromFill(fill model.FillEvent) error {
	return nil
}

func (p *portfolio) isInvested(symbol string) (model.Position, bool) {
	// Todo: Test this func asap rocky
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