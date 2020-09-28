package storage

import (
	"errors"
	"time"
)

var (
	ErrUserExists        = errors.New("user already exists")
	ErrUserDoesNotExist  = errors.New("user doesn't not exist")
	ErrEmptyEvents       = errors.New("empty events list")
	ErrEventDoesNotExist = errors.New("event doesn't exist")
)

type Event struct {
	ID       int64     `db:"id"`
	UserID   int64     `db:"user_id"`
	Title    string    `db:"title"`
	Content  string    `db:"content"`
	DateFrom time.Time `db:"date_from"`
	DateTo   time.Time `db:"date_to"`
	Notified bool      `db:"notified"`
}

// email must be unique and not null.
type User struct {
	ID        int64  `db:"id"`
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

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
