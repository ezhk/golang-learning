package internalhttp

import (
	"fmt"
	"net/http"

	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/viper"
)

type HTTPServer struct {
	Server *http.Server
}

func NewHTTPServer(log *logger.Logger) *HTTPServer {
	handler := &ServeHandler{}

	mux := http.NewServeMux()

	helloHandler := http.HandlerFunc(handler.Hello)
	mux.Handle("/hello", LoggerMiddleware(log, helloHandler))

	address := fmt.Sprintf("%s:%d",
		viper.GetString("server.host"),
		viper.GetInt("server.port"))
	server := &http.Server{Addr: address, Handler: mux}

	return &HTTPServer{Server: server}
}

func (s *HTTPServer) Run() error {
	return s.Server.ListenAndServe()
}
