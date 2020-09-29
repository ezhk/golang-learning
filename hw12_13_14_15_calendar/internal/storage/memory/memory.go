package memorystorage

import (
	"sync"
	"time"

	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	utils "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/utils"
)

type MemoryDatabase struct {
	Events map[int64]*storage.Event
	Users  map[int64]*storage.User

	EventsByUserID map[int64][]*storage.Event
	UsersByEmail   map[string]*storage.User

	LatestRecordID int64
	LatestUserID   int64

	mutex sync.RWMutex
}

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
	m.mutex.RLock()
	userPtr, ok := m.UsersByEmail[email]
	m.mutex.RUnlock()

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

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.LatestUserID++
	user.ID = m.LatestUserID

	m.Users[user.ID] = &user
	m.UsersByEmail[user.Email] = &user

	return user, nil
}

func (m *MemoryDatabase) UpdateUser(user storage.User) error {
	m.mutex.RLock()
	userPointer, ok := m.Users[user.ID]
	m.mutex.RUnlock()

	if !ok {
		return storage.ErrUserDoesNotExist
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.UsersByEmail, userPointer.Email)

	m.UsersByEmail[user.Email] = &user
	m.Users[user.ID] = &user

	return nil
}

func (m *MemoryDatabase) DeleteUser(user storage.User) error {
	m.mutex.RLock()
	userPointer, ok := m.Users[user.ID]
	m.mutex.RUnlock()

	if !ok {
		return storage.ErrUserDoesNotExist
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.UsersByEmail, userPointer.Email)
	delete(m.Users, userPointer.ID)

	return nil
}

func (m *MemoryDatabase) GetEventsByUserID(userID int64) ([]storage.Event, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

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
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	events, ok := m.EventsByUserID[userID]
	if !ok {
		return nil, storage.ErrEmptyEvents
	}

	selectedEvents := make([]storage.Event, 0)
	for _, eventPointer := range events {
		if eventPointer.DateFrom.After(fromDate) && eventPointer.DateFrom.Before(toDate) {
			selectedEvents = append(selectedEvents, *eventPointer)

			continue
		}

		if eventPointer.DateTo.After(fromDate) && eventPointer.DateTo.Before(toDate) {
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

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.LatestRecordID++
	event.ID = m.LatestRecordID

	m.Events[event.ID] = &event
	m.EventsByUserID[userID] = append(m.EventsByUserID[userID], &event)

	return event, nil
}

func (m *MemoryDatabase) UpdateEvent(event storage.Event) error {
	m.mutex.RLock()
	eventPointer, ok := m.Events[event.ID]
	m.mutex.RUnlock()

	if !ok {
		return storage.ErrEventDoesNotExist
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

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
	m.mutex.RLock()
	eventPointer, ok := m.Events[event.ID]
	m.mutex.RUnlock()

	if !ok {
		return storage.ErrEventDoesNotExist
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for idx, pointers := range m.EventsByUserID[event.UserID] {
		if pointers != eventPointer {
			continue
		}
		m.EventsByUserID[event.UserID] = append(m.EventsByUserID[event.UserID][:idx], m.EventsByUserID[event.UserID][idx+1:]...)
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
	m.mutex.RLock()
	defer m.mutex.RUnlock()

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
	m.mutex.RLock()
	eventPointer, ok := m.Events[event.ID]
	m.mutex.RUnlock()

	if !ok {
		return storage.ErrEventDoesNotExist
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Back compatibility with SQL.
	event.Notified = true
	eventPointer.Notified = true

	return nil
}
