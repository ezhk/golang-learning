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
	query := `SELECT *
	FROM users
	WHERE email = $1`

	var user storage.User
	err := d.database.QueryRowxContext(d.ctx, query, email).StructScan(&user)

	switch {
	case err == sql.ErrNoRows:
		return storage.User{}, storage.ErrUserDoesNotExist
	case err != nil:
		return storage.User{}, fmt.Errorf("query error: %q", err)
	}

	return user, nil
}

func (d *SQLDatabase) CreateUser(email string, firstName string, lastName string) (int64, error) {
	existUser, _ := d.GetUserByEmail(email)
	if v := existUser.ID; v > 0 {
		return v, storage.ErrUserExists
	}

	var insertID int64
	query := `INSERT INTO users(
		email,
		first_name,
		last_name
	) VALUES (
		$1, $2, $3
	) RETURNING id`

	err := d.database.QueryRowContext(d.ctx, query, email, firstName, lastName).Scan(&insertID)
	if err != nil {
		return -1, fmt.Errorf("create user error: %q", err)
	}

	return insertID, nil
}

func (d *SQLDatabase) UpdateUser(id int64, email string, firstName string, lastName string) error {
	query := `UPDATE users
	SET	email = :email,
		first_name = :first_name,
		last_name = :last_name
	WHERE id = :id`

	u := storage.User{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
	_, err := d.database.NamedExecContext(d.ctx, query, &u)
	if err != nil {
		return fmt.Errorf("update user error: %q", err)
	}

	return nil
}

func (d *SQLDatabase) DeleteUserByUserID(id int64) error {
	query := `DELETE FROM users
	WHERE id = $1`

	_, err := d.database.ExecContext(d.ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user error: %q", err)
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
		return nil, err
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
		return nil, err
	}

	return events, nil
}

func (d *SQLDatabase) CreateEvent(userID int64, title, content string, dateFrom, dateTo time.Time) (int64, error) {
	var insertID int64
	query := `INSERT INTO events(
		user_id,
		title,
		content,
		date_from,
		date_to
	) VALUES (
		$1, $2, $3, $4, $5
	) RETURNING id`

	err := d.database.QueryRowContext(d.ctx, query, userID, title, content, dateFrom, dateTo).Scan(&insertID)
	if err != nil {
		return -1, fmt.Errorf("create event error: %q", err)
	}

	return insertID, nil
}

func (d *SQLDatabase) UpdateEvent(id int64, userID int64, title, content string, dateFrom, dateTo time.Time) error {
	query := `UPDATE events
	SET	user_id = :user_id,
		title = :title,
		content = :content,
		date_from = :date_from,
		date_to = :date_to
	WHERE id = :id`

	e := storage.Event{
		ID:       id,
		UserID:   userID,
		Title:    title,
		Content:  content,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	_, err := d.database.NamedExecContext(d.ctx, query, &e)
	if err != nil {
		return fmt.Errorf("update event error: %q", err)
	}

	return nil
}

func (d *SQLDatabase) DeleteEvent(id int64) error {
	query := `DELETE
	FROM events
	WHERE id = $1`

	_, err := d.database.ExecContext(d.ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete event error: %q", err)
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
		return nil, err
	}

	return events, nil
}

func (d *SQLDatabase) MarkEventAsNotified(id int64) error {
	query := `UPDATE events
	SET notified = TRUE
	WHERE id = $1`

	_, err := d.database.ExecContext(d.ctx, query, id)
	if err != nil {
		return fmt.Errorf("update event as notified error: %q", err)
	}

	return nil
}
