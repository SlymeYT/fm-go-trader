package portfolio

import (
	"github.com/pkg/errors"
	"github.com/sheerun/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/data"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/model"
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

func (p *portfolio) GenerateOrders() error {
	return nil
}

func (p *portfolio) UpdateFromFill() error {
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

func NewPortfolio() *portfolio {
	return &portfolio{}
}

func parseAdvise(signalPairs map[string]float32) string {
	return "LONG" // LONG, CLOSE_LONG, SHORT or CLOSE_SHORT
}