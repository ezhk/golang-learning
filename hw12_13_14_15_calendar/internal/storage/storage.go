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
	ID       int64     `db:"id" json:"id"`
	UserID   int64     `db:"user_id" json:"user_id"`
	Title    string    `db:"title" json:"title"`
	Content  string    `db:"content" json:"content"`
	DateFrom time.Time `db:"date_from" json:"date_from"`
	DateTo   time.Time `db:"date_to" json:"date_to"`
	Notified bool      `db:"notified" json:"notified"`
}

// email must be unique and not null.
type User struct {
	ID        int64  `db:"id" json:"id"`
	Email     string `db:"email" json:"email"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
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
