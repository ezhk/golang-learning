package internalgrpc

import (
	structs "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/structs"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertInternalUserToUser(u *User) structs.User {
	return structs.User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func ConvertUserToInternalUser(u structs.User) *User {
	return &User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func CovertInternalEventToEvent(e *Event) structs.Event {
	return structs.Event{
		ID:       e.ID,
		UserID:   e.UserID,
		Title:    e.Title,
		Content:  e.Content,
		DateFrom: e.DateFrom.AsTime(),
		DateTo:   e.DateTo.AsTime(),
		Notified: e.Notified,
	}
}

func ConvertEventToInternalEvent(e structs.Event) *Event {
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
