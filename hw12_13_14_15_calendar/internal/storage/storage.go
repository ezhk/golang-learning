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
	CreateUser(string, string, string) (int64, error)
	UpdateUser(int64, string, string, string) error
	DeleteUser(int64) error

	GetRecordsByUserID(int64) ([]Event, error)
	CreateRecord(int64, string, string, time.Time, time.Time) (int64, error)
	UpdateRecord(int64, int64, string, string, time.Time, time.Time) error
	DeleteRecord(int64) error
}
