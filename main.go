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

	fmt.Printf("%+v", cfg)

	return nil
}
