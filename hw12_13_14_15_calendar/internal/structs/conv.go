package structs

import (
	"fmt"
	"strconv"
	"time"
)

// ConvertUserMessageToUser — parse UserMessage (from JSON) to User.
func ConvertUserMessageToUser(u UserMessage) (User, error) {
	ID, err := strconv.Atoi(u.ID)
	if err != nil {
		return User{}, fmt.Errorf("convert string to int error: %w", err)
	}

	user := User{
		ID:        int64(ID),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}

	return user, nil
}

// ConvertEventMessageToEvent — parse EventMessage (from JSON) to Event.
func ConvertEventMessageToEvent(e EventMessage) (Event, error) {
	ID, err := strconv.Atoi(e.ID)
	if err != nil {
		return Event{}, fmt.Errorf("convert ID error: %w", err)
	}

	userID, err := strconv.Atoi(e.ID)
	if err != nil {
		return Event{}, fmt.Errorf("convert UserID error: %w", err)
	}

	dateFrom, err := time.Parse(time.RFC3339, e.DateFrom)
	if err != nil {
		return Event{}, fmt.Errorf("parse DateFrom error: %w", err)
	}

	dateTo, err := time.Parse(time.RFC3339, e.DateTo)
	if err != nil {
		return Event{}, fmt.Errorf("parse DateTo error: %w", err)
	}

	event := Event{
		ID:       int64(ID),
		UserID:   int64(userID),
		Title:    e.Title,
		Content:  e.Content,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Notified: e.Notified,
	}

	return event, nil
}
