package database

import (
	"errors"
	"time"
)

type MemoryDatabase struct {
	Calendar
	Users []User

	LatestRecordID int
	LatestUserID   int
}

func NewDatatabase() ClientInterface {
	return &MemoryDatabase{}
}

func (m *MemoryDatabase) Connect(_ string) error {
	return nil
}

func (m *MemoryDatabase) Close() error {
	return nil
}

func (m *MemoryDatabase) GetUser(firstName string, lastName string) User {
	for _, u := range m.Users {
		if u.FirstName == firstName && u.LastName == lastName {
			return u
		}
	}

	return User{}
}

func (m *MemoryDatabase) CreateUser(firstName string, lastName string) (int, error) {
	existUser := m.GetUser(firstName, lastName)
	if v := existUser.ID; v > 0 {
		return v, errors.New("user already exists")
	}

	m.LatestUserID++
	m.Users = append(m.Users, User{ID: m.LatestUserID, FirstName: firstName, LastName: lastName})

	return m.LatestUserID, nil
}

func (m *MemoryDatabase) UpdateUser(id int, firstName string, lastName string) error {
	for idx, u := range m.Users {
		if u.ID == id {
			m.Users[idx] = User{
				ID:        id,
				FirstName: firstName,
				LastName:  lastName,
			}

			return nil
		}
	}

	return errors.New("user not found")
}

func (m *MemoryDatabase) DeleteUser(id int) error {
	for idx, u := range m.Users {
		if u.ID == id {
			m.Users = append(m.Users[:idx], m.Users[idx+1:]...)

			return nil
		}
	}

	return errors.New("user not found")
}

func (m *MemoryDatabase) GetRecords(userID int) []CRecord {
	findRecords := make([]CRecord, 0)

	for _, r := range m.Records {
		if r.userID == userID {
			findRecords = append(findRecords, r)
		}
	}

	return findRecords
}

func (m *MemoryDatabase) CreateRecord(userID int, title, content string) (int, error) {
	m.LatestRecordID++
	m.Records = append(m.Records, CRecord{
		ID:        m.LatestRecordID,
		userID:    userID,
		Title:     title,
		Content:   content,
		UpdatedAt: time.Now(),
	})

	return m.LatestRecordID, nil
}

func (m *MemoryDatabase) UpdateRecord(id int, title, content string) error {
	for idx, r := range m.Records {
		if r.ID == id {
			m.Records[idx] = CRecord{
				ID:        id,
				userID:    m.Records[idx].userID,
				Title:     title,
				Content:   content,
				UpdatedAt: time.Now(),
			}

			return nil
		}
	}

	return errors.New("record not found")
}

func (m *MemoryDatabase) DeleteRecord(id int) error {
	for idx, r := range m.Records {
		if r.ID == id {
			m.Records = append(m.Records[:idx], m.Records[idx+1:]...)

			return nil
		}
	}

	return errors.New("record not found")
}
