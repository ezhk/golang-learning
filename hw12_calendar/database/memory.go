package database

import (
	"errors"
	"time"
)

type MemoryDatabase struct {
	Calendar
	Users []User

	LatestRecordId int
	LatestUserId   int
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
	if v := existUser.Id; v > 0 {
		return v, errors.New("user already exists")
	}

	m.LatestUserId++
	m.Users = append(m.Users, User{Id: m.LatestUserId, FirstName: firstName, LastName: lastName})

	return m.LatestUserId, nil
}

func (m *MemoryDatabase) UpdateUser(id int, firstName string, lastName string) error {
	for idx, u := range m.Users {
		if u.Id == id {
			m.Users[idx] = User{
				Id:        id,
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
		if u.Id == id {
			m.Users = append(m.Users[:idx], m.Users[idx+1:]...)
			return nil
		}
	}

	return errors.New("user not found")
}

func (m *MemoryDatabase) GetRecords(userId int) []CRecord {
	findRecords := make([]CRecord, 0)

	for _, r := range m.Records {
		if r.UserId == userId {
			findRecords = append(findRecords, r)
		}
	}

	return findRecords
}

func (m *MemoryDatabase) CreateRecord(userId int, title, content string) (int, error) {
	m.LatestRecordId++
	m.Records = append(m.Records, CRecord{
		Id:        m.LatestRecordId,
		UserId:    userId,
		Title:     title,
		Content:   content,
		UpdatedAt: time.Now(),
	})

	return m.LatestRecordId, nil
}

func (m *MemoryDatabase) UpdateRecord(id int, title, content string) error {
	for idx, r := range m.Records {
		if r.Id == id {
			m.Records[idx] = CRecord{
				Id:        id,
				UserId:    m.Records[idx].UserId,
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
		if r.Id == id {
			m.Records = append(m.Records[:idx], m.Records[idx+1:]...)
			return nil
		}
	}

	return errors.New("record not found")
}
