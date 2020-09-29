package internalhttphandlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
)

type OneEventBody struct {
	Status string        `json:"status"`
	Event  storage.Event `json:"event"`
}

type ManyEventsBody struct {
	Status string          `json:"status"`
	Events []storage.Event `json:"events"`
}

func generateOneEvent(w io.Writer, e storage.Event) error {
	return json.NewEncoder(w).Encode(OneEventBody{
		Status: StatusOK,
		Event:  e,
	})
}

func generateManyEvent(w io.Writer, e []storage.Event) error {
	return json.NewEncoder(w).Encode(ManyEventsBody{
		Status: StatusOK,
		Events: e,
	})
}

func (sh *ServeHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	event, err = sh.db.CreateEvent(event.UserID, event.Title, event.Content, event.DateFrom, event.DateTo)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	_ = generateOneEvent(w, event)
}

func (sh *ServeHandler) ReadEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	v := r.URL.Query()

	filter := v.Get("filter")
	userID, _ := strconv.ParseInt(v.Get("userid"), 10, 64)

	var events []storage.Event
	var err error

	switch filter {
	case "daily":
		events, err = sh.db.DailyEvents(userID, time.Now())
	case "weekly":
		events, err = sh.db.WeeklyEvents(userID, time.Now())
	case "monthly":
		events, err = sh.db.MonthlyEvents(userID, time.Now())
	case "uninformed":
		events, err = sh.db.GetNotifyReadyEvents()
	default:
		events, err = sh.db.GetEventsByUserID(userID)
	}

	if err != nil {
		_ = generateError(w, err)

		return
	}

	_ = generateManyEvent(w, events)
}

func (sh *ServeHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var event storage.Event

	params := mux.Vars(r)
	eventID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	event.ID = eventID
	err = sh.db.UpdateEvent(event)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	_ = generateOneEvent(w, event)
}

func (sh *ServeHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	eventID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	err = sh.db.DeleteEvent(storage.Event{ID: eventID})
	if err != nil {
		_ = generateError(w, err)

		return
	}

	_ = generateOneEvent(w, storage.Event{ID: eventID})
}
