package internalgrpc

import (
	context "context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/server/http"
	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate protoc -I . -I ${GOPATH}/src -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. --go_out=plugins=grpc:. --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true internalgrpc.proto

type Server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener

	cfg *config.Configuration
	log *logger.Logger

	db storage.ClientInterface
}

func NewServer(cfg *config.Configuration, log *logger.Logger, db storage.ClientInterface) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:    ctx,
		cancel: cancel,

		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (s *Server) RunProxy() error {
	mux := runtime.NewServeMux()

	err := RegisterCalendarHandlerServer(s.ctx, mux, s)
	if err != nil {
		return err
	}
	address := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.HTTPPort)

	// Endble HTTP logger.
	return http.ListenAndServe(address, internalhttp.LoggerMiddleware(s.log, mux))
}

func (s *Server) RunServer() error {
	address := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.GRPCPort)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil
	}

	s.listener = l

	// Include gRPC logger.
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(&s.log.Logger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(&s.log.Logger),
		)),
	)
	RegisterCalendarServer(grpcServer, s)

	return grpcServer.Serve(s.listener)
}

func (s *Server) Close() error {
	s.cancel()

	return s.listener.Close()
}

func covertToStorageUser(u *User) storage.User {
	return storage.User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func convertToServerUser(u storage.User) *User {
	return &User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func (s *Server) GetUser(ctx context.Context, req *RequestByUserEmail) (*User, error) {
	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	return convertToServerUser(user), nil
}

func (s *Server) CreateUser(ctx context.Context, u *User) (*User, error) {
	user, err := s.db.CreateUser(u.Email, u.FirstName, u.LastName)
	if err != nil {
		return nil, err
	}

	return convertToServerUser(user), nil
}

func (s *Server) UpdateUser(ctx context.Context, u *User) (*User, error) {
	user := covertToStorageUser(u)

	err := s.db.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *RequestByUserID) (*User, error) {
	err := s.db.DeleteUser(storage.User{ID: req.ID})
	if err != nil {
		return nil, err
	}

	return &User{ID: req.ID}, nil
}

func covertToStorageEvent(e *Event) storage.Event {
	return storage.Event{
		ID:       e.ID,
		UserID:   e.UserID,
		Title:    e.Title,
		Content:  e.Content,
		DateFrom: e.DateFrom.AsTime(),
		DateTo:   e.DateTo.AsTime(),
		Notified: e.Notified,
	}
}

func convertToServerEvent(e storage.Event) *Event {
	return &Event{
		ID:       e.ID,
		UserID:   e.UserID,
		Title:    e.Title,
		Content:  e.Content,
		DateFrom: timestamppb.New(e.DateFrom),
		DateTo:   timestamppb.New(e.DateTo),
		Notified: e.Notified,
	}
}

func (s *Server) GetEvents(ctx context.Context, req *RequestByUserID) (*Events, error) {
	events, err := s.db.GetEventsByUserID(req.ID)
	if err != nil {
		return nil, err
	}

	resultEvents := make([]*Event, 0)
	for _, e := range events {
		resultEvents = append(resultEvents, convertToServerEvent(e))
	}

	return &Events{Events: resultEvents}, nil
}

func (s *Server) CreateEvent(ctx context.Context, e *Event) (*Event, error) {
	event, err := s.db.CreateEvent(e.UserID, e.Title, e.Content, e.DateFrom.AsTime(), e.DateTo.AsTime())
	if err != nil {
		return nil, err
	}

	return convertToServerEvent(event), nil
}

func (s *Server) UpdateEvent(ctx context.Context, e *Event) (*Event, error) {
	err := s.db.UpdateEvent(covertToStorageEvent(e))
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *Server) DeleteEvent(ctx context.Context, eventID *EventID) (*Event, error) {
	event := storage.Event{ID: eventID.ID}
	err := s.db.DeleteEvent(event)
	if err != nil {
		return nil, err
	}

	return convertToServerEvent(event), nil
}

func (s *Server) PeriodEvents(ctx context.Context, d *DateEvent) (*Events, error) {
	var events []storage.Event
	var err error

	switch d.Period {
	case DateEvent_DAILY:
		events, err = s.db.DailyEvents(d.UserID, d.Date.AsTime())
	case DateEvent_WEEKLY:
		events, err = s.db.DailyEvents(d.UserID, d.Date.AsTime())
	case DateEvent_MONTHLY:
		events, err = s.db.DailyEvents(d.UserID, d.Date.AsTime())
	default:
		return nil, errors.New("not period filter")
	}

	if err != nil {
		return nil, err
	}

	resultEvents := make([]*Event, 0)
	for _, e := range events {
		resultEvents = append(resultEvents, convertToServerEvent(e))
	}

	return &Events{Events: resultEvents}, nil
}

func (s *Server) GetNotifyReadyEvents(ctx context.Context, empty *emptypb.Empty) (*Events, error) {
	events, err := s.db.GetNotifyReadyEvents()
	if err != nil {
		return nil, err
	}

	resultEvents := make([]*Event, 0)
	for _, e := range events {
		resultEvents = append(resultEvents, convertToServerEvent(e))
	}

	return &Events{Events: resultEvents}, nil
}

// func (s *Server) MarkEventAsNotified(ctx context.Context, e *Event) (*emptypb.Empty, error) {
// 	event := covertToStorageEvent(e)
// 	err := s.db.MarkEventAsNotified(&event)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return nil, nil
// }
