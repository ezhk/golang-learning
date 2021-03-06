package internalgrpc

import (
	context "context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	config "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/server/http"
	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	structs "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/structs"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate protoc -I . -I ${GOPATH}/src -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. --go_out=plugins=grpc:. --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true internalgrpc.proto

type (
	ErrorMessage      = structs.ErrorMessage
	UserMessage       = structs.UserMessage
	EventMessage      = structs.EventMessage
	ManyEventsMessage = structs.ManyEventsMessage
)

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
	// Enable omitted fields to output.
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			OrigName:     true,
			EmitDefaults: true,
		}))

	err := RegisterCalendarHandlerServer(s.ctx, mux, s)
	if err != nil {
		return fmt.Errorf("register handler error: %w", err)
	}
	address := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.HTTPPort)

	// Endble HTTP logger.
	return http.ListenAndServe(address, internalhttp.LoggerMiddleware(s.log, mux))
}

func (s *Server) RunServer() error {
	address := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.GRPCPort)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
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

func (s *Server) GetUser(ctx context.Context, req *RequestByUserEmail) (*User, error) {
	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("get user error: %w", err)
	}

	return ConvertUserToInternalUser(user), nil
}

func (s *Server) CreateUser(ctx context.Context, u *User) (*User, error) {
	user, err := s.db.CreateUser(u.Email, u.FirstName, u.LastName)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}

	return ConvertUserToInternalUser(user), nil
}

func (s *Server) UpdateUser(ctx context.Context, u *User) (*User, error) {
	user := ConvertInternalUserToUser(u)

	err := s.db.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("update user error: %w", err)
	}

	return u, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *RequestByUserID) (*User, error) {
	err := s.db.DeleteUser(structs.User{ID: req.ID})
	if err != nil {
		return nil, fmt.Errorf("delete user error: %w", err)
	}

	return &User{ID: req.ID}, nil
}

func (s *Server) GetEvents(ctx context.Context, req *RequestByUserID) (*Events, error) {
	events, err := s.db.GetEventsByUserID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("get events error: %w", err)
	}

	resultEvents := make([]*Event, 0)
	for _, e := range events {
		resultEvents = append(resultEvents, ConvertEventToInternalEvent(e))
	}

	return &Events{Events: resultEvents}, nil
}

func (s *Server) CreateEvent(ctx context.Context, e *Event) (*Event, error) {
	event, err := s.db.CreateEvent(e.UserID, e.Title, e.Content, e.DateFrom.AsTime(), e.DateTo.AsTime())
	if err != nil {
		return nil, fmt.Errorf("create event error: %w", err)
	}

	return ConvertEventToInternalEvent(event), nil
}

func (s *Server) UpdateEvent(ctx context.Context, e *Event) (*Event, error) {
	err := s.db.UpdateEvent(CovertInternalEventToEvent(e))
	if err != nil {
		return nil, fmt.Errorf("update event error: %w", err)
	}

	return e, nil
}

func (s *Server) DeleteEvent(ctx context.Context, eventID *EventID) (*Event, error) {
	event := structs.Event{ID: eventID.ID}
	err := s.db.DeleteEvent(event)
	if err != nil {
		return nil, fmt.Errorf("delete event error: %w", err)
	}

	return ConvertEventToInternalEvent(event), nil
}

func (s *Server) PeriodEvents(ctx context.Context, d *DateEvent) (*Events, error) {
	var events []structs.Event
	var err error

	dateTime, err := time.Parse(time.RFC3339, d.Date)
	if err != nil {
		return nil, fmt.Errorf("parse time error: %w", err)
	}

	switch d.Period {
	case DateEvent_DAILY:
		events, err = s.db.DailyEvents(d.UserID, dateTime)
	case DateEvent_WEEKLY:
		events, err = s.db.WeeklyEvents(d.UserID, dateTime)
	case DateEvent_MONTHLY:
		events, err = s.db.MonthlyEvents(d.UserID, dateTime)
	default:
		return nil, errors.New("not period filter")
	}

	if err != nil {
		return nil, fmt.Errorf("select period events error: %w", err)
	}

	resultEvents := make([]*Event, 0)
	for _, e := range events {
		resultEvents = append(resultEvents, ConvertEventToInternalEvent(e))
	}

	return &Events{Events: resultEvents}, nil
}

func (s *Server) GetNotifyReadyEvents(ctx context.Context, empty *emptypb.Empty) (*Events, error) {
	events, err := s.db.GetNotifyReadyEvents()
	if err != nil {
		return nil, fmt.Errorf("get notify ready events error: %w", err)
	}

	resultEvents := make([]*Event, 0)
	for _, e := range events {
		resultEvents = append(resultEvents, ConvertEventToInternalEvent(e))
	}

	return &Events{Events: resultEvents}, nil
}
