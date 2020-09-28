package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	config "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
)

type HTTPServer struct {
	Server *http.Server
}

func NewHTTPServer(cfg *config.Configuration, log *logger.Logger) *HTTPServer {
	handler := &ServeHandler{}

	mux := http.NewServeMux()

	helloHandler := http.HandlerFunc(handler.Hello)
	mux.Handle("/hello", LoggerMiddleware(log, helloHandler))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: mux,
	}

	return &HTTPServer{Server: server}
}

func (s *HTTPServer) Run() error {
	return s.Server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
