package api

type apiError struct {
	Error 	string	`json:"Error"`
	Message string	`json:"Message"`
}

//// server is a HTTP server.
//type server struct {
//	http.Server
//	router  *chi.Mux
//	log     *zap.Logger
//	engine  service.TradingEngine
//	name    string
//	version string
//}
//
//// NewServer constructs a new server.
//func NewServer(cfg *config.Server, log *zap.Logger, engine service.TradingEngine) *server {
//	s := &server{
//		log:     log,
//		engine:  engine,
//		name:    cfg.Name,
//		version: cfg.Version,
//	}
//	s.router = chi.NewRouter()
//	s.Server = http.Server{
//		Addr: fmt.Sprintf(":%d", cfg.Port),
//		Handler: s,
//	}
//	s.registerRoutes()
//
//	return s
//}
//
//func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	s.router.ServeHTTP(w, r)
//}
//
//func (s *server) Run() {
//	// Todo: Handle graceful shutdown with channel pattern!
//	s.log.Info(fmt.Sprintf("%s-%s running on port %s", s.name, s.version, s.Addr))
//	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//		s.log.Error(err.Error())
//	}
//}
//
//func (s *server) respondError(w http.ResponseWriter, code int, err error, message string) {
//	s.log.Error(err.Error())
//	s.respondJSON(w, code, apiError{
//		Error: err.Error(),
//		Message: message,
//	})
//}
//
//func (s *server) respondJSON(w http.ResponseWriter, code int, payload interface{}) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(code)
//	if err := json.NewEncoder(w).Encode(payload); err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		_, _ = w.Write([]byte(err.Error()))
//		return
//	}
//}