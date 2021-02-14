package portfolio

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sheerun/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/config"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/data"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/strategy"
	"go.uber.org/zap"
	"math"
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
	data              data.Handler
	sizeManager       SizeManager
	riskManager       RiskManager
	symbol            string
	initialCash       float64
	currentCash       float64
	currentValue      float64
	orders            []model.OrderEvent
	fills             []model.FillEvent
	positions         map[string]model.Position
	historicPositions map[string][]model.Position
}

// UpdateFromMarket updates the current portfolio positions using the new market event data
func (p *portfolio) UpdateFromMarket(market model.MarketEvent) error {
	// Update current positions
	if position, isInvested := p.isInvested(p.symbol); isInvested {
		err := position.Update(market)
		if err != nil {
			return errors.Wrap(err, "failed portfolio.UpdateFromMarket()")
		}
		p.positions[p.symbol] = position
	}

	return nil
}

// isInvested determines if a portfolio has an open Position for a Symbol & returns that position
func (p *portfolio) isInvested(symbol string) (model.Position, bool) {
	// Todo: Test this func asap rocky
	position, isInPositions := p.positions[symbol]
	// If present in current positions & exit fill value is zero
	if isInPositions && position.ExitFillValueNet == 0 {
		return position, true
	}
	return position, false
}

// GenerateOrders parses a SignalEvent and generates an OrderEvent if the portfolio wants to act on the signal advise
func (p *portfolio) GenerateOrders(signal model.SignalEvent) error {
	// Todo: Enhance this to allow for closing a trade and opening a reverse trade on the same market event

	// Check if the SignalEvent is for a Symbol already invested in
	position, isInvested := p.isInvested(signal.Symbol)

	// If no cash, cannot open a new position -> exit without generating an order
	if !isInvested && p.currentCash == 0.0 {
		return nil
	}

	// Parse SignalPairs map to determine the net OrderEvent decision
	strength, decision := p.parseSignalDecisions(position, isInvested, signal.SignalPairs)

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

// parseSignalDecisions assesses what, if any, decisions should be made based on incoming signalPairs
func (p *portfolio) parseSignalDecisions(position model.Position, isInvested bool, signalPairs map[string]float32) (float32, string) {
	// Todo: Test this func asap rocky

	// Pull (strength, decisionAdvise) out of signalPairs map
	strengthLong, long := signalPairs[strategy.DecisionLong]
	strengthCloseLong, closeLong := signalPairs[strategy.DecisionCloseLong]
	strengthShort, short := signalPairs[strategy.DecisionShort]
	strengthCloseShort, closeShort := signalPairs[strategy.DecisionCloseShort]

	if isInvested && position.Direction == "LONG" {
		if closeLong {
			return strengthCloseLong, strategy.DecisionCloseLong
		}
	}

	if isInvested && position.Direction == "SHORT" {
		if closeShort {
			return strengthCloseShort, strategy.DecisionCloseShort
		}
	}

	if !isInvested {
		if long {
			return strengthLong, strategy.DecisionLong
		} else if short {
			return strengthShort, strategy.DecisionShort
		}
	}

	return 0.0, strategy.DecisionNothing
}

// UpdateFromFill updates the portfolio's current positions & historicPositions from a FillEvent
func (p *portfolio) UpdateFromFill(fill model.FillEvent) error {

	// Get current available data to determine the FillValueGross - would be determined in execution for live trading
	currentData, latestBarIndex := p.data.GetLatestData()
	fill.FillValueGross = math.Abs(fill.Quantity) * currentData.Closes[latestBarIndex]

	// Must be an exit
	if position, isInvested := p.isInvested(fill.Symbol); isInvested {
		// Exit position instance
		err := position.Exit(fill)
		if err != nil {
			return errors.Wrap(err, "failed portfolio.UpdateFromFill()")
		}
		// Append exited position to historicPositions and remove from current positions
		p.historicPositions[fill.Symbol] = append(p.historicPositions[fill.Symbol], position)
		delete(p.positions, fill.Symbol)

		// Update cash & value on exit
		p.currentCash = p.currentCash + fill.FillValueGross
		p.currentValue = p.currentCash

	} else {
		// Must be an entry
		position := model.Position{}
		err := position.Enter(fill)
		if err != nil {
			return errors.Wrap(err, "failed portfolio.UpdateFromFill()")
		}

		// Update cash & value on entry
		p.currentCash = p.currentCash - fill.FillValueGross
		p.currentValue = p.currentCash + fill.FillValueGross
	}

	// Update completed FillEvents
	p.fills = append(p.fills, fill)

	p.log.Info(fmt.Sprintf("Value: %v, Cash: %v, Holdings: %+v", p.currentValue, p.currentCash, p.positions))

	return nil
}

func NewPortfolio(cfg config.Trader, eventQ *queue.Queue, data data.Handler) *portfolio {
	return &portfolio{
		log:              cfg.Log,
		eventQ:           eventQ,
		data:              data,
		sizeManager:       &Size{DefaultOrderValue: cfg.DefaultOrderValue},
		riskManager:       &Risk{DefaultOrderType: OrderTypeMarket},
		symbol:            cfg.Symbol,
		initialCash:       cfg.StartingCash,
		currentCash:       cfg.StartingCash,
		currentValue:      cfg.StartingCash,
		orders:            []model.OrderEvent{},
		fills:             []model.FillEvent{},
		positions:         make(map[string]model.Position),
		historicPositions: make(map[string][]model.Position),
	}
}