package sqlstorage

/*
	Docker postgres run:
	docker run --name golang-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres:alpine
	docker exec -ti golang-postgres psql -U postgres
*/

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	utils "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/utils"

	// use pgx as a database/sql compatible driver.
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type SQLDatabase struct {
	database *sqlx.DB

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
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}
	d.database = db

	return nil
}

func (d *SQLDatabase) Close() error {
	d.cancel()

	return d.database.Close()
}

func (d *SQLDatabase) GetUserByEmail(email string) (storage.User, error) {
	query := `SELECT *
	FROM users
	WHERE email = $1`

	var user storage.User
	err := d.database.QueryRowxContext(d.ctx, query, email).StructScan(&user)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return storage.User{}, storage.ErrUserDoesNotExist
	case err != nil:
		return storage.User{}, fmt.Errorf("get user by email error: %w", err)
	}

	return user, nil
}

func (d *SQLDatabase) CreateUser(email string, firstName string, lastName string) (storage.User, error) {
	existUser, _ := d.GetUserByEmail(email)
	if v := existUser.ID; v > 0 {
		return existUser, storage.ErrUserExists
	}

	user := storage.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	var insertID int64
	query := `INSERT INTO users(
		email,
		first_name,
		last_name
	) VALUES (
		:email,
		:first_name,
		:last_name
	) RETURNING id`

	// Prepare named request.
	nstmt, err := d.database.PrepareNamedContext(d.ctx, query)
	if err != nil {
		return storage.User{}, fmt.Errorf("prepare create user error: %w", err)
	}

	// Execute and waiting for returning ID.
	err = nstmt.Get(&insertID, &user)
	if err != nil {
		return storage.User{}, fmt.Errorf("create user error: %w", err)
	}

	user.ID = insertID

	return user, nil
}

func (d *SQLDatabase) UpdateUser(user storage.User) error {
	query := `UPDATE users
	SET	email = :email,
		first_name = :first_name,
		last_name = :last_name
	WHERE id = :id`

	_, err := d.database.NamedExecContext(d.ctx, query, &user)
	if err != nil {
		return fmt.Errorf("update user error: %w", err)
	}

	return nil
}

func (d *SQLDatabase) DeleteUser(user storage.User) error {
	query := `DELETE FROM users
	WHERE id = $1`

	_, err := d.database.ExecContext(d.ctx, query, user.ID)
	if err != nil {
		return fmt.Errorf("delete user error: %w", err)
	}

	return nil
}

func (d *SQLDatabase) GetEventsByUserID(userID int64) ([]storage.Event, error) {
	query := `SELECT *
	FROM events
	WHERE user_id = $1`

	var events []storage.Event
	err := d.database.SelectContext(d.ctx, &events, query, userID)
	if err != nil {
		return nil, fmt.Errorf("select events by user ID error: %w", err)
	}

	return events, nil
}

func (d *SQLDatabase) getEventsByDateRange(userID int64, fromDate, toDate time.Time) ([]storage.Event, error) {
	query := `SELECT *
	FROM events
	WHERE user_id = $1
		AND (
			date_from >= $2 AND date_from < $3
			OR date_to >= $2 AND date_to < $3
		)`

	var events []storage.Event
	err := d.database.SelectContext(d.ctx, &events, query, userID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("select events by date range error: %w", err)
	}

	return events, nil
}

func (d *SQLDatabase) CreateEvent(userID int64, title, content string, dateFrom, dateTo time.Time) (storage.Event, error) {
	event := storage.Event{
		UserID:   userID,
		Title:    title,
		Content:  content,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	var insertID int64
	query := `INSERT INTO events(
		user_id,
		title,
		content,
		date_from,
		date_to
	) VALUES (
		:user_id,
		:title,
		:content,
		:date_from,
		:date_to
	) RETURNING id`

	// Prepare named request.
	nstmt, err := d.database.PrepareNamedContext(d.ctx, query)
	if err != nil {
		return storage.Event{}, fmt.Errorf("prepare create event error: %w", err)
	}

	// Execute and waiting for returning ID.
	err = nstmt.Get(&insertID, &event)
	if err != nil {
		return storage.Event{}, fmt.Errorf("create event error: %w", err)
	}

	event.ID = insertID

	return event, nil
}

func (d *SQLDatabase) UpdateEvent(event storage.Event) error {
	query := `UPDATE events
	SET	user_id = :user_id,
		title = :title,
		content = :content,
		date_from = :date_from,
		date_to = :date_to,
		notified = :notified
	WHERE id = :id`

	_, err := d.database.NamedExecContext(d.ctx, query, &event)
	if err != nil {
		return fmt.Errorf("update event error: %w", err)
	}

	return nil
}

func (d *SQLDatabase) DeleteEvent(event storage.Event) error {
	query := `DELETE
	FROM events
	WHERE id = $1`

	_, err := d.database.ExecContext(d.ctx, query, event.ID)
	if err != nil {
		return fmt.Errorf("delete event error: %w", err)
	}

	return nil
}

func (d *SQLDatabase) DailyEvents(userID int64, date time.Time) ([]storage.Event, error) {
	fromDate := utils.StartDay(date)
	toDate := utils.EndDay(date)

	return d.getEventsByDateRange(userID, fromDate, toDate)
}

func (d *SQLDatabase) WeeklyEvents(userID int64, date time.Time) ([]storage.Event, error) {
	fromDate := utils.StartWeek(date)
	toDate := utils.EndWeek(date)

	return d.getEventsByDateRange(userID, fromDate, toDate)
}

func (d *SQLDatabase) MonthlyEvents(userID int64, date time.Time) ([]storage.Event, error) {
	fromDate := utils.StartMonth(date)
	toDate := utils.EndMonth(date)

	return d.getEventsByDateRange(userID, fromDate, toDate)
}

func (d *SQLDatabase) GetNotifyReadyEvents() ([]storage.Event, error) {
	query := `SELECT *
	FROM events
	WHERE notified = FALSE
		AND date_from <= $1`

	twoWeekLeft := time.Now().Add(2 * time.Hour * 24 * 7)

	var events []storage.Event
	err := d.database.SelectContext(d.ctx, &events, query, twoWeekLeft)
	if err != nil {
		return nil, fmt.Errorf("get notify ready events error: %w", err)
	}

	return events, nil
}

func (d *SQLDatabase) MarkEventAsNotified(event *storage.Event) error {
	query := `UPDATE events
	SET notified = TRUE
	WHERE id = $1`

	_, err := d.database.ExecContext(d.ctx, query, event.ID)
	if err != nil {
		return fmt.Errorf("update events error: %w", err)
	}
	event.Notified = true

	return nil
}
