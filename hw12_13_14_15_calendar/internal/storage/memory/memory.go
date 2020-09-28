package memorystorage

import (
	"sync"
	"time"

	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
)

type MemoryDatabase struct {
	Events map[int64]storage.Event
	Users  map[int64]storage.User

	EventsByUserID map[int64][]int64
	UsersIDByEmail map[string]int64

	LatestRecordID int64
	LatestUserID   int64

	mutex sync.RWMutex
}

func NewDatatabase() storage.ClientInterface {
	return &MemoryDatabase{
		Events:         make(map[int64]storage.Event),
		Users:          make(map[int64]storage.User),
		EventsByUserID: make(map[int64][]int64),
		UsersIDByEmail: make(map[string]int64),
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
	userID, ok := m.UsersIDByEmail[email]
	m.mutex.RUnlock()

	if !ok {
		return storage.User{}, storage.ErrUserDoesNotExist
	}

	m.mutex.RLock()
	user, ok := m.Users[userID]
	m.mutex.RUnlock()

	if ok {
		return user, nil
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.UsersIDByEmail, email)

	return storage.User{}, storage.ErrUserDoesNotExist
}

func (m *MemoryDatabase) CreateUser(email string, firstName string, lastName string) (int64, error) {
	existUser, err := m.GetUserByEmail(email)
	if err == nil {
		return existUser.ID, storage.ErrUserExists
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.LatestUserID++
	m.Users[m.LatestUserID] = storage.User{
		ID:        m.LatestUserID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
	m.UsersIDByEmail[email] = m.LatestUserID

	return m.LatestUserID, nil
}

func (m *MemoryDatabase) UpdateUser(id int64, email string, firstName string, lastName string) error {
	m.mutex.RLock()
	user, ok := m.Users[id]
	m.mutex.RUnlock()

	if !ok {
		return storage.ErrUserDoesNotExist
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.UsersIDByEmail, user.Email)
	m.UsersIDByEmail[email] = id
	m.Users[id] = storage.User{
		ID:        m.LatestUserID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	return nil
}

func (m *MemoryDatabase) DeleteUser(id int64) error {
	m.mutex.RLock()
	user, ok := m.Users[id]
	m.mutex.RUnlock()

	if !ok {
		return storage.ErrUserDoesNotExist
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.UsersIDByEmail, user.Email)
	delete(m.Users, id)

	return nil
}

func (m *MemoryDatabase) GetRecordsByUserID(userID int64) ([]storage.Event, error) {
	m.mutex.RLock()
	events, ok := m.EventsByUserID[userID]
	m.mutex.RUnlock()

	if !ok {
		return nil, storage.ErrEmptyEvents
	}

	selectedEvents := make([]storage.Event, 0)
	for idx, eventID := range events {
		m.mutex.RLock()
		event, ok := m.Events[eventID]
		m.mutex.RUnlock()

		if !ok {
			m.mutex.Lock()
			m.EventsByUserID[userID] = append(m.EventsByUserID[userID][:idx], m.EventsByUserID[userID][idx+1:]...)
			m.mutex.Unlock()

			continue
		}

		selectedEvents = append(selectedEvents, event)
	}

	return selectedEvents, nil
}

func (m *MemoryDatabase) CreateRecord(userID int64, title, content string, dateFrom, dateTo time.Time) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.LatestRecordID++
	m.Events[m.LatestRecordID] = storage.Event{
		ID:       m.LatestRecordID,
		UserID:   userID,
		Title:    title,
		Content:  content,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	m.EventsByUserID[userID] = append(m.EventsByUserID[userID], m.LatestRecordID)

	return m.LatestRecordID, nil
}

func (m *MemoryDatabase) UpdateRecord(id int64, userID int64, title, content string, dateFrom, dateTo time.Time) error {
	m.mutex.RLock()
	event, ok := m.Events[id]
	m.mutex.RUnlock()

	if !ok {
		return storage.ErrEventDoesNotExist
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.EventsByUserID, event.UserID)
	m.EventsByUserID[userID] = append(m.EventsByUserID[userID], id)
	m.Events[id] = storage.Event{
		ID:       id,
		UserID:   userID,
		Title:    title,
		Content:  content,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	return nil
}

func (m *MemoryDatabase) DeleteRecord(id int64) error {
	m.mutex.RLock()
	event, ok := m.Events[id]
	m.mutex.RUnlock()

	if !ok {
		return storage.ErrEventDoesNotExist
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for idx, eventID := range m.EventsByUserID[event.UserID] {
		if eventID != id {
			continue
		}
		m.EventsByUserID[event.UserID] = append(m.EventsByUserID[event.UserID][:idx], m.EventsByUserID[event.UserID][idx+1:]...)
	}

	delete(m.Events, id)

	return nil
}
