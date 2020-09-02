package database

import "time"

type Calendar struct {
	Records []CRecord
}

type CRecord struct {
	ID        int       `db:"id"`
	userID    int       `db:"user_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	UpdatedAt time.Time `db:"updated_at"`
}

type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

type ClientInterface interface {
	Connect(string) error
	Close() error

	GetUser(string, string) User
	CreateUser(string, string) (int, error)
	UpdateUser(int, string, string) error
	DeleteUser(int) error

	GetRecords(int) []CRecord
	CreateRecord(int, string, string) (int, error)
	UpdateRecord(int, string, string) error
	DeleteRecord(int) error
}
