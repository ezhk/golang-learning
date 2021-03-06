package memorystorage

import (
	"time"

	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	utils "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/utils"
)

type MemoryDatabase storage.MemoryDatabase

func NewDatatabase() storage.ClientInterface {
	return &MemoryDatabase{
		Events:         make(map[int64]*storage.Event),
		Users:          make(map[int64]*storage.User),
		EventsByUserID: make(map[int64][]*storage.Event),
		UsersByEmail:   make(map[string]*storage.User),
	}
}

func (m *MemoryDatabase) Connect(_ string) error {
	return nil
}

func (m *MemoryDatabase) Close() error {
	return nil
}

func (m *MemoryDatabase) GetUserByEmail(email string) (storage.User, error) {
	m.RLock()
	userPtr, ok := m.UsersByEmail[email]
	m.RUnlock()

	if !ok {
		return storage.User{}, storage.ErrUserDoesNotExist
	}

	return *userPtr, nil
}

func (m *MemoryDatabase) CreateUser(email string, firstName string, lastName string) (storage.User, error) {
	existUser, err := m.GetUserByEmail(email)
	if err == nil {
		return existUser, storage.ErrUserExists
	}

	user := storage.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	m.Lock()
	defer m.Unlock()

	m.LatestUserID++
	user.ID = m.LatestUserID

	m.Users[user.ID] = &user
	m.UsersByEmail[user.Email] = &user

	return user, nil
}

func (m *MemoryDatabase) UpdateUser(user storage.User) error {
	m.RLock()
	userPointer, ok := m.Users[user.ID]
	m.RUnlock()

	if !ok {
		return storage.ErrUserDoesNotExist
	}

	m.Lock()
	defer m.Unlock()

	delete(m.UsersByEmail, userPointer.Email)

	m.UsersByEmail[user.Email] = &user
	m.Users[user.ID] = &user

	return nil
}

func (m *MemoryDatabase) DeleteUser(user storage.User) error {
	m.RLock()
	userPointer, ok := m.Users[user.ID]
	m.RUnlock()

	if !ok {
		return storage.ErrUserDoesNotExist
	}

	m.Lock()
	defer m.Unlock()

	delete(m.UsersByEmail, userPointer.Email)
	delete(m.Users, userPointer.ID)

	return nil
}

func (m *MemoryDatabase) GetEventsByUserID(userID int64) ([]storage.Event, error) {
	m.RLock()
	defer m.RUnlock()

	events, ok := m.EventsByUserID[userID]
	if !ok {
		return nil, storage.ErrEmptyEvents
	}

	selectedEvents := make([]storage.Event, 0)
	for _, eventPointer := range events {
		selectedEvents = append(selectedEvents, *eventPointer)
	}

	return selectedEvents, nil
}

func (m *MemoryDatabase) getEventsByDateRange(userID int64, fromDate, toDate time.Time) ([]storage.Event, error) {
	m.RLock()
	defer m.RUnlock()

	events, ok := m.EventsByUserID[userID]
	if !ok {
		return nil, storage.ErrEmptyEvents
	}

	selectedEvents := make([]storage.Event, 0)
	for _, eventPointer := range events {
		if eventPointer.DateFrom.Before(toDate) && eventPointer.DateTo.After(fromDate) {
			selectedEvents = append(selectedEvents, *eventPointer)

			continue
		}
	}

	return selectedEvents, nil
}

func (m *MemoryDatabase) CreateEvent(userID int64, title, content string, dateFrom, dateTo time.Time) (storage.Event, error) {
	event := storage.Event{
		UserID:   userID,
		Title:    title,
		Content:  content,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	m.Lock()
	defer m.Unlock()

	m.LatestRecordID++
	event.ID = m.LatestRecordID

	m.Events[event.ID] = &event
	m.EventsByUserID[userID] = append(m.EventsByUserID[userID], &event)

	return event, nil
}

func (m *MemoryDatabase) UpdateEvent(event storage.Event) error {
	m.RLock()
	eventPointer, ok := m.Events[event.ID]
	m.RUnlock()

	if !ok {
		return storage.ErrEventDoesNotExist
	}

	m.Lock()
	defer m.Unlock()

	// Change attributes by exist pointer.
	eventPointer.UserID = event.UserID
	eventPointer.Title = event.Title
	eventPointer.Content = event.Content
	eventPointer.DateFrom = event.DateFrom
	eventPointer.DateTo = event.DateTo
	eventPointer.Notified = event.Notified

	return nil
}

func (m *MemoryDatabase) DeleteEvent(event storage.Event) error {
	m.RLock()
	eventPointer, ok := m.Events[event.ID]
	m.RUnlock()

	if !ok {
		return storage.ErrEventDoesNotExist
	}

	m.Lock()
	defer m.Unlock()

	for idx, pointers := range m.EventsByUserID[eventPointer.UserID] {
		if pointers != eventPointer {
			continue
		}

		m.EventsByUserID[eventPointer.UserID] = append(
			m.EventsByUserID[eventPointer.UserID][:idx],
			m.EventsByUserID[eventPointer.UserID][idx+1:]...,
		)

		break
	}

	delete(m.Events, event.ID)

	return nil
}

func (m *MemoryDatabase) DailyEvents(userID int64, date time.Time) ([]storage.Event, error) {
	fromDate := utils.StartDay(date)
	toDate := utils.EndDay(date)

	return m.getEventsByDateRange(userID, fromDate, toDate)
}

func (m *MemoryDatabase) WeeklyEvents(userID int64, date time.Time) ([]storage.Event, error) {
	fromDate := utils.StartWeek(date)
	toDate := utils.EndWeek(date)

	return m.getEventsByDateRange(userID, fromDate, toDate)
}

func (m *MemoryDatabase) MonthlyEvents(userID int64, date time.Time) ([]storage.Event, error) {
	fromDate := utils.StartMonth(date)
	toDate := utils.EndMonth(date)

	return m.getEventsByDateRange(userID, fromDate, toDate)
}

func (m *MemoryDatabase) GetNotifyReadyEvents() ([]storage.Event, error) {
	m.RLock()
	defer m.RUnlock()

	twoWeekLeft := time.Now().Add(2 * time.Hour * 24 * 7)
	selectedEvents := make([]storage.Event, 0)

	for _, eventPointer := range m.Events {
		if !eventPointer.Notified && eventPointer.DateFrom.Before(twoWeekLeft) {
			selectedEvents = append(selectedEvents, *eventPointer)
		}
	}

	return selectedEvents, nil
}

func (m *MemoryDatabase) MarkEventAsNotified(event *storage.Event) error {
	m.RLock()
	eventPointer, ok := m.Events[event.ID]
	m.RUnlock()

	if !ok {
		return storage.ErrEventDoesNotExist
	}

	m.Lock()
	defer m.Unlock()

	// Back compatibility with SQL.
	event.Notified = true
	eventPointer.Notified = true

	return nil
}
