package structs

import (
	"sync"
	"time"
)

// Event - database event struct.
type Event struct {
	ID       int64     `db:"id" json:"ID"`
	UserID   int64     `db:"user_id" json:"UserID"`
	Title    string    `db:"title" json:"Title"`
	Content  string    `db:"content" json:"Content"`
	DateFrom time.Time `db:"date_from" json:"DateFrom"`
	DateTo   time.Time `db:"date_to" json:"DateTo"`
	Notified bool      `db:"notified" json:"Notified"`
}

// User is storage user struct.
// email must be unique and not null.
type User struct {
	ID        int64  `db:"id" json:"ID"`
	Email     string `db:"email" json:"Email"`
	FirstName string `db:"first_name" json:"FirstName"`
	LastName  string `db:"last_name" json:"LastName"`
}

// MemoryDatabase - in memory database base struct.
type MemoryDatabase struct {
	sync.RWMutex

	Events map[int64]*Event
	Users  map[int64]*User

	EventsByUserID map[int64][]*Event
	UsersByEmail   map[string]*User

	LatestRecordID int64
	LatestUserID   int64
}
