package api

import "net/http"

func (s *server) runBacktest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Running... Backtest"))
		//s.engine.RunBacktest()
	}
}

func (s *server) runTraderLive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Run Trader Live"))
	}
}

func (s *server) runTraderDry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Run Trader Dry"))
	}
}
