package main

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/config"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	log, err := zap.NewDevelopment()
	if err != nil {
		return errors.Wrap(err, "failed to init logger")
	}
	defer log.Sync()

	cfg, err := config.GetConfig(log)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to init environment config: %s", err))
	}

	//traderService, err := service.NewTradingEngine(&cfg.Engine, log)
	//if err != nil {
	//	log.Fatal(fmt.Sprintf("failed to init trading engine: %s", err))
	//}
	//
	//server := api.NewServer(&cfg.Server, log, traderService)
	//server.Run()

	return nil
}

