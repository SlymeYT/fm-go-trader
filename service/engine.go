package service

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/config"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/trader"
	"go.uber.org/zap"
)

type TradingEngine interface {
	RunBacktest() error
	RunTraderLive() error
	RunTraderDry() error
}

type tradingEngine struct {
	log     *zap.Logger
	traders []trader.Trader
}

func (t *tradingEngine) RunBacktest() error {
	for _, traderPair := range t.traders {
		if err := traderPair.Run(); err != nil {
			return errors.Wrap(err, "failed to RunBacktest()")
		}
	}
	return nil
}

func (t *tradingEngine) RunTraderLive() error {
	// Run live trading with meta-strategy for specified trading pairs
	return nil
}

func (t *tradingEngine) RunTraderDry() error {
	// Run live trading with meta-strategy for specified trading pairs
	return nil
}

func NewTradingEngine(cfg *config.Engine, log *zap.Logger) (*tradingEngine, error) {
	traders, err := buildTraders(cfg, log)
	if err != nil {
		return &tradingEngine{}, err
	}

	engine := &tradingEngine{
		log:     log,
		traders: traders,
	}

	return engine, nil
}

func buildTraders(cfg *config.Engine, log *zap.Logger) ([]trader.Trader, error) {
	var traders []trader.Trader
	for index, cfg := range buildTraderConfigs(cfg, log) {
		traderPair, err := trader.NewTrader(cfg)
		if err != nil {
			return traders, errors.Wrap(err, fmt.Sprintf("failed to init trader %v with config: %+v\n", index, cfg))
		}
		traders = append(traders, traderPair)
	}
	return traders, nil
}

func buildTraderConfigs(cfg *config.Engine, log *zap.Logger) []config.Trader {
	// Todo: Remove these stubs with the real array from the config.Engine struct
	tickers := []string{cfg.Symbols}
	timeframes := []string{cfg.Timeframes}
	exchanges := []string{cfg.Exchanges}
	startingCash := cfg.StartingCash / float64(len(tickers))
	defaultOrderValue := startingCash / 10

	var traderConfigs []config.Trader
	for index, symbol := range tickers {
		traderConfigs = append(traderConfigs, config.Trader{
			Log: 				log,
			Symbol:    			symbol,
			Timeframe:      	timeframes[index],
			Exchange:       	exchanges[index],
			StartingCash: 		startingCash,
			DefaultOrderValue: 	defaultOrderValue,
		})
	}
	return traderConfigs
}