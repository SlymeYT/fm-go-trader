package strategy

import (
	"github.com/sheerun/queue"
	"gitlab.com/open-source-keir/financial-modelling/trading/fm-trader/data"
	"go.uber.org/zap"
)

type Strategy interface {
	GenerateSignal() error
}

type simpleRSIStrategy struct {
	log          *zap.Logger
	eventQ       *queue.Queue
	data         *data.Handler
}

