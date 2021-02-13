package api

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"time"
)

func (s *server) registerRoutes() {
	apiVersion := fmt.Sprintf("/api/%s", s.version)

	// Middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))  	// request timout

	// Command
	s.router.Route(apiVersion + "/backtest", func(r chi.Router) {
		r.Get("/", s.runBacktest())  				// POST /api/v1/backtest
	})

	s.router.Route(apiVersion + "/trader", func(r chi.Router) {
		r.Get("/live", s.runTraderLive())  				// POST /api/v1/trader/live
		r.Get("/dry", s.runTraderDry())  				// POST /api/v1/trader/dry
	})

	// Query
}