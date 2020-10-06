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
	ID       int64     `db:"id" json:"ID"`
	UserID   int64     `db:"user_id" json:"UserID"`
	Title    string    `db:"title" json:"Title"`
	Content  string    `db:"content" json:"Content"`
	DateFrom time.Time `db:"date_from" json:"DateFrom"`
	DateTo   time.Time `db:"date_to" json:"DateTo"`
	Notified bool      `db:"notified" json:"Notified"`
}

// email must be unique and not null.
type User struct {
	ID        int64  `db:"id" json:"ID"`
	Email     string `db:"email" json:"Email"`
	FirstName string `db:"first_name" json:"FirstName"`
	LastName  string `db:"last_name" json:"LastName"`
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
