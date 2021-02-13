package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	ActiveProfile = "ACTIVE_PROFILE"
	Directory     = "config"
)

// Config is the complete configuration loaded from the file associated with ActiveProfile.
type Config struct {
	Engine Engine
	Server Server
}

// config.Repository is the repository & persistence configuration.
type Repository struct {
}

// config.Engine is the engine service configuration.
type Engine struct {
	// Tickers is the array of tickers the trading engine will use to create Traders
	Symbols string			`envconfig:"TICKERS" required:"true"`
	// Timeframes is the array of timeframe the trading engine will use to create Traders
	Timeframes string		`envconfig:"TIMEFRAMES" required:"true"`
	// Exchanges is the array of exchanges the trading engine will use to create Traders
	Exchanges string 		`envconfig:"EXCHANGES" required:"true"`
	// Starting
	StartingCapital float64 `envconfig:"STARTING_CAPITAL" required:"true"`
}

type Trader struct {
	// Ticker is the ticker symbol this instance of Trader is using.
	Symbol string
	// Timeframe is the interval between bars this instance of Trader is using.
	Timeframe string
	// Exchange is the name of the exchange this instance of Trader is using.
	Exchange string
	// StartingCapital is the starting capital allocated to this instance of Trader.
	StartingCash float64
}

// config.Server is the HTTP server configuration.
type Server struct {
	// Name is the Name of the function the Server is hosting.
	Name string 		`envconfig:"SERVER_NAME" required:"true"`
	// Version is the Version of the servic
	Version string 		`envconfig:"SERVER_VERSION" required:"true"`
	// Port is the HTTP Port to serve on
	Port int 			`envconfig:"SERVER_PORT" required:"true"`
}

func GetConfig(log *zap.Logger) (*Config, error) {
	activeProfile := strings.TrimSpace(os.Getenv(ActiveProfile))

	if activeProfile == "" {
		activeProfile = "default"
	}

	file := filepath.Join(Directory, fmt.Sprintf("%s.env", activeProfile))
	if err := godotenv.Load(file); err != nil {
		return nil, err
	}

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return &config, err
	}

	if err := validateConfig(config); err != nil {
		return &config, errors.Wrap(err, "failed validation")
	}

	log.Info(fmt.Sprintf("environment configuration loaded: %s", file))
	return &config, nil
}

func validateConfig(cfg Config) error {
	refl := reflect.ValueOf(cfg)
	for i := 0; i < refl.NumField(); i++ {
		for j := 0; j < refl.Field(i).NumField(); j++ {
			if refl.Field(i).Field(j).Interface() == "" {
				return errors.New(fmt.Sprintf("config field %s cannot be empty", refl.Field(i).Type().Field(j).Name))
			}
		}
	}
	return nil
}


