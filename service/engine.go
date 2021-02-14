package service

import (
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/trader"
	"go.uber.org/zap"
)

type TradingEngine interface {
	RunBacktest()
	RunTraderLive()
	RunTraderDry()
}

type tradingEngine struct {
	log     *zap.Logger
	traders []trader.Trader
}

func (t *tradingEngine) RunBacktest() {
	// Run backtest of meta-strategy for specified trading pairs (on go routine with unique id that can be cancelled)
	// 1. Parse meta-strategy from static config or http request
	// 2. Parse symbols to trade from static config or http request
	// 3. Spin up pair Traders on go routines
	// 4. Report results

	for _, traderPair := range t.traders {
		traderPair.Run()
	}
}