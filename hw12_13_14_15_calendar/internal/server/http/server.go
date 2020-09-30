package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	config "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
	internalhttphandlers "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/server/http/handlers"
	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
)

/*
REST API scheme and storage methods.

/users			- GET		=> 404 Not Found
/users/{email}	- GET		=> GetUserByEmail(string) (User, error)
/users			- POST		=> CreateUser(string, string, string) (User, error)
/users/{id}		- PUT		=> UpdateUser(User) error
/users/{id}		- DELETE	=> DeleteUser(User{ID: id}) error


/events	- GET	=> 404 Not Found

/events?userid=ID					- GET	=> GetEventsByUserID(int64) ([]Event, error)
/events?userid=ID&filter=daily?date=2006-01-02T12:04:05Z03:00	- GET	=> DailyEvents(userID int64, date time.Time) ([]Event, error)
/events?userid=ID&filter=weekly?date=2006-01-02T12:04:05Z03:00	- GET	=> WeeklyEvents(userID int64, date time.Time) ([]Event, error)
/events?userid=ID&filter=monthly?date=2006-01-02T12:04:05Z03:00	- GET	=> MonthlyEvents(userID int64, date time.Time) ([]Event, error)

/events?filter=uninformed	— GET	=> GetNotifyReadyEvents() ([]Event, error)

/events			- POST		=> CreateEvent(int64, string, string, time.Time, time.Time) (Event, error)
/events/{id}	- PUT		=> UpdateEvent(Event) error
/events/{id}	- DELETE	=> DeleteEvent(Event{ID: id}) error
*/

type HTTPServer struct {
	Server *http.Server
}

func NewHTTPServer(cfg *config.Configuration, log *logger.Logger, db storage.ClientInterface) *HTTPServer {
	handler := internalhttphandlers.NewServeHandler(log, db)

	r := mux.NewRouter()
	// Users methods.
	r.HandleFunc("/users", handler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{email}", handler.ReadUser).Methods("GET")
	r.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")
	// Events methods.
	r.HandleFunc("/events", handler.CreateEvent).Methods("POST")
	r.HandleFunc("/events", handler.ReadEvents).Methods("GET")
	r.HandleFunc("/events/{id}", handler.UpdateEvent).Methods("PUT")
	r.HandleFunc("/events/{id}", handler.DeleteEvent).Methods("DELETE")

	helloHandler := http.HandlerFunc(handler.HealthCheck)

	mux := http.NewServeMux()
	mux.Handle("/health-check", LoggerMiddleware(log, helloHandler))
	mux.Handle("/", LoggerMiddleware(log, r))

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
