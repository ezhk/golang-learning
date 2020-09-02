package server

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type HTTPServer struct {
	Server *http.Server
}

func NewHTTPServer() *HTTPServer {
	handler := &ServeHandler{}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handler.Hello)

	address := fmt.Sprintf("%s:%d",
		viper.GetString("server.host"),
		viper.GetInt("server.port"))
	server := &http.Server{Addr: address, Handler: mux}

	return &HTTPServer{Server: server}
}

func (s *HTTPServer) Run() error {
	return s.Server.ListenAndServe()
}
