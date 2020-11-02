package server

import (
	"context"
	"net"
	"net/http"

	"github.com/ezhk/golang-learning/banners-rotation/internal/api"
	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/logger"
	"github.com/ezhk/golang-learning/banners-rotation/internal/queue"
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
	storage storage.DatabaseInterface
	queue   *queue.Queue
}

func NewServer(configPtr *config.Configuration, loggerPtr *logger.Logger, databaseI storage.DatabaseInterface, queuePrt *queue.Queue) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:    ctx,
		cancel: cancel,

		config:  configPtr,
		logger:  loggerPtr,
		storage: databaseI,
		queue:   queuePrt,
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
	api.RegisterBannerServer(grpcServer, s)

	return grpcServer.Serve(s.listener)
}

func (s Server) RunProxy() error {
	mux := runtime.NewServeMux()

	err := api.RegisterBannerHandlerServer(s.ctx, mux, s)
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
