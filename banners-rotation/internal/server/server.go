package server

import (
	"context"
	"net"
	"net/http"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/logger"
	"github.com/ezhk/golang-learning/banners-rotation/internal/storage"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener

	config  *config.Configuration
	logger  *logger.Logger
	storage *storage.Storage
}

//go:generate protoc -I . -I ${GOPATH}/src -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go-grpc_out . --go-grpc_opt require_unimplemented_servers=false --go_out . --go_opt paths=source_relative --openapiv2_out . --openapiv2_opt logtostderr=true --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true server.proto

func NewServer(configPtr *config.Configuration, loggerPtr *logger.Logger, storagePtr *storage.Storage) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:    ctx,
		cancel: cancel,

		config:  configPtr,
		logger:  loggerPtr,
		storage: storagePtr,
	}
}

func (s Server) Run() error {
	l, err := net.Listen("tcp", s.config.Server.GRPC)
	if err != nil {
		return err
	}
	s.listener = l

	// Include gRPC logger.
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(&s.logger.Logger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(&s.logger.Logger),
		)),
	)
	RegisterBannerServer(grpcServer, s)

	return grpcServer.Serve(s.listener)
}

func (s Server) RunProxy() error {
	mux := runtime.NewServeMux()

	err := RegisterBannerHandlerServer(s.ctx, mux, s)
	if err != nil {
		return err
	}

	// Endble HTTP logger.
	return http.ListenAndServe(s.config.Server.HTTP, ProxyLoggerMiddleware(s.logger, mux))
}

func (s Server) Close() error {
	s.cancel()

	return s.listener.Close()
}
