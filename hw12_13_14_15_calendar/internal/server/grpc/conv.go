package internalgrpc

import (
	"strconv"
	"time"

	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUserMessageToStorageUser(u UserMessage) (storage.User, error) {
	ID, err := strconv.Atoi(u.ID)
	if err != nil {
		return storage.User{}, err
	}

	storageUser := storage.User{
		ID:        int64(ID),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}

	return storageUser, nil
}

func ConvertEventMessageToStorageUser(e EventMessage) (storage.Event, error) {
	ID, err := strconv.Atoi(e.ID)
	if err != nil {
		return storage.Event{}, err
	}

	userID, err := strconv.Atoi(e.ID)
	if err != nil {
		return storage.Event{}, err
	}

	dateFrom, err := time.Parse(time.RFC3339, e.DateFrom)
	if err != nil {
		return storage.Event{}, err
	}

	dateTo, err := time.Parse(time.RFC3339, e.DateTo)
	if err != nil {
		return storage.Event{}, err
	}

	storageEvent := storage.Event{
		ID:       int64(ID),
		UserID:   int64(userID),
		Title:    e.Title,
		Content:  e.Content,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Notified: e.Notified,
	}

	return storageEvent, nil
}

func CovertUserToStorageUser(u *User) storage.User {
	return storage.User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func ConvertStorageUserToUser(u storage.User) *User {
	return &User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func CovertEventToStorageEvent(e *Event) storage.Event {
	return storage.Event{
		ID:       e.ID,
		UserID:   e.UserID,
		Title:    e.Title,
		Content:  e.Content,
		DateFrom: e.DateFrom.AsTime(),
		DateTo:   e.DateTo.AsTime(),
		Notified: e.Notified,
	}
}

func ConvertStorageEventToEvent(e storage.Event) *Event {
	return &Event{
		ID:       e.ID,
		UserID:   e.UserID,
		Title:    e.Title,
		Content:  e.Content,
		DateFrom: timestamppb.New(e.DateFrom),
		DateTo:   timestamppb.New(e.DateTo),
		Notified: e.Notified,
	}
}
