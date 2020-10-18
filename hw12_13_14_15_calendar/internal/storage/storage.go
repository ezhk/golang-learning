package storage

import (
	"errors"
	"time"

	structs "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/structs"
)

var (
	ErrUserExists        = errors.New("user already exists")
	ErrUserDoesNotExist  = errors.New("user doesn't not exist")
	ErrEmptyEvents       = errors.New("empty events list")
	ErrEventDoesNotExist = errors.New("event doesn't exist")
)

type (
	Event          = structs.Event
	User           = structs.User
	MemoryDatabase = structs.MemoryDatabase
)

type ClientInterface interface {
	Connect(string) error
	Close() error

	GetUserByEmail(string) (User, error)
	CreateUser(string, string, string) (User, error)
	UpdateUser(User) error
	DeleteUser(User) error

	GetEventsByUserID(int64) ([]Event, error)
	CreateEvent(int64, string, string, time.Time, time.Time) (Event, error)
	UpdateEvent(Event) error
	DeleteEvent(Event) error

	DailyEvents(userID int64, date time.Time) ([]Event, error)
	WeeklyEvents(userID int64, date time.Time) ([]Event, error)
	MonthlyEvents(userID int64, date time.Time) ([]Event, error)

	GetNotifyReadyEvents() ([]Event, error)
	MarkEventAsNotified(event *Event) error
}
