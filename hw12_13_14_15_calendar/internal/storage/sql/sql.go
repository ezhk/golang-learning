package sqlstorage

/*
	Docker postgres run:
	docker run --name golang-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres:alpine
	docker exec -ti golang-postgres psql -U postgres
*/

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"

	// use pgx as a database/sql compatible driver.
	_ "github.com/jackc/pgx/stdlib"
)

type SQLDatabase struct {
	database *sql.DB

	ctx    context.Context
	cancel context.CancelFunc
}

func NewDatatabase() storage.ClientInterface {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &SQLDatabase{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (d *SQLDatabase) Connect(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to load driver: %q", err)
	}
	d.database = db

	return nil
}

func (d *SQLDatabase) Close() error {
	d.cancel()

	return d.database.Close()
}

func (d *SQLDatabase) GetUserByEmail(email string) (storage.User, error) {
	query := "SELECT * FROM users WHERE email = $1"

	var (
		id                  int64
		firstName, lastName string
	)

	err := d.database.QueryRowContext(d.ctx, query, email).Scan(&id, &email, &firstName, &lastName)
	switch {
	case err == sql.ErrNoRows:
		return storage.User{}, storage.ErrUserDoesNotExist
	case err != nil:
		return storage.User{}, fmt.Errorf("query error: %q", err)
	}

	return storage.User{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}

func (d *SQLDatabase) CreateUser(email string, firstName string, lastName string) (int64, error) {
	existUser, _ := d.GetUserByEmail(email)
	if v := existUser.ID; v > 0 {
		return v, storage.ErrUserExists
	}

	var insertID int64
	query := "INSERT INTO users(email, first_name, last_name) VALUES ($1, $2, $3) RETURNING id"

	err := d.database.QueryRowContext(d.ctx, query, email, firstName, lastName).Scan(&insertID)
	if err != nil {
		return -1, fmt.Errorf("create user error: %q", err)
	}

	return insertID, nil
}

func (d *SQLDatabase) UpdateUser(id int64, email string, firstName string, lastName string) error {
	query := "UPDATE users SET email = $1, first_name = $2, last_name = $3 WHERE id = $4"
	_, err := d.database.ExecContext(d.ctx, query, email, firstName, lastName, id)
	if err != nil {
		return fmt.Errorf("update user error: %q", err)
	}

	return nil
}

func (d *SQLDatabase) DeleteUser(id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := d.database.ExecContext(d.ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user error: %q", err)
	}

	return nil
}

func (d *SQLDatabase) GetRecordsByUserID(userID int64) ([]storage.Event, error) {
	query := "SELECT * FROM events WHERE user_id = $1"
	rows, err := d.database.QueryContext(d.ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	findRecords := make([]storage.Event, 0)
	for rows.Next() {
		var (
			id, userID       int64
			title, content   string
			dateFrom, dateTo time.Time
		)

		err = rows.Scan(&id, &userID, &title, &content, &dateFrom, &dateTo)
		if err != nil {
			return findRecords, err
		}

		findRecords = append(findRecords, storage.Event{
			ID:       id,
			UserID:   userID,
			Title:    title,
			Content:  content,
			DateFrom: dateFrom,
			DateTo:   dateTo,
		})
	}

	if rows.Err() != nil {
		return findRecords, rows.Err()
	}

	return findRecords, nil
}

func (d *SQLDatabase) CreateRecord(userID int64, title, content string, dateFrom, dateTo time.Time) (int64, error) {
	var insertID int64
	query := "INSERT INTO events(user_id, title, content, date_from, date_to) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err := d.database.QueryRowContext(d.ctx, query, userID, title, content, dateFrom, dateTo).Scan(&insertID)
	if err != nil {
		return -1, fmt.Errorf("create event error: %q", err)
	}

	return insertID, nil
}

func (d *SQLDatabase) UpdateRecord(id int64, userID int64, title, content string, dateFrom, dateTo time.Time) error {
	query := "UPDATE events SET user_id = $1, title = $2, content = $3, date_from = $4, date_to = $5 WHERE id = $6"
	_, err := d.database.ExecContext(d.ctx, query, userID, title, content, dateFrom, dateTo, id)
	if err != nil {
		return fmt.Errorf("update event error: %q", err)
	}

	return nil
}

func (d *SQLDatabase) DeleteRecord(id int64) error {
	query := "DELETE FROM events WHERE id = $1"
	_, err := d.database.ExecContext(d.ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete event error: %q", err)
	}

	return nil
}
