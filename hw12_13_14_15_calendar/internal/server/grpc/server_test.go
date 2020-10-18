package internalgrpc

import (
	"bytes"
	context "context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	config "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	structs "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/structs"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	mux  *runtime.ServeMux
	user structs.User
}

func TestUserSuite(t *testing.T) {
	// This is what actually runs our suite
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) SetupTest() {
	cfg := config.NewConfig("testdata/config.yaml")

	// Init in-memory database.
	database := cfg.DatabaseBuilder()
	err := database.Connect(cfg.GetDatabasePath())
	s.NoError(err)

	// Default defined user.
	s.user = structs.User{
		Email:     "test@yandex.ru",
		FirstName: "Test",
		LastName:  "Yandex mail",
	}

	// Init server.
	ctx, cancel := context.WithCancel(context.Background())
	s.mux = runtime.NewServeMux()
	srv := &Server{
		ctx:    ctx,
		cancel: cancel,

		db: database,
	}

	// Register server.
	RegisterCalendarHandlerServer(context.TODO(), s.mux, srv)
}

func (s *UserTestSuite) TestEmptyUser() {
	r := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	s.mux.ServeHTTP(w, r)

	s.Equal(http.StatusMethodNotAllowed, w.Result().StatusCode)
}

func (s *UserTestSuite) TestCreateUser() {
	// Create user.
	jsonBody, err := json.Marshal(s.user)
	s.NoError(err)
	r, err := http.NewRequest("POST", "/api/v1/users", bytes.NewReader(jsonBody))
	s.NoError(err)
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Result().StatusCode)

	var userMsg UserMessage
	err = json.NewDecoder(w.Body).Decode(&userMsg)
	s.NoError(err)
	user, err := structs.ConvertUserMessageToUser(userMsg)
	s.NoError(err)

	s.Equal(int64(1), user.ID)
	s.Equal(s.user.Email, user.Email)
	s.Equal(s.user.FirstName, user.FirstName)
	s.Equal(s.user.LastName, user.LastName)

	// Create user again might cause error.
	r, err = http.NewRequest("POST", "/api/v1/users", bytes.NewReader(jsonBody))
	s.NoError(err)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	s.Equal(http.StatusInternalServerError, w.Result().StatusCode)
	var errMsg ErrorMessage
	err = json.NewDecoder(w.Body).Decode(&errMsg)
	s.NoError(err)
	s.Equal("user already exists", errMsg.Message)
	s.Equal(2, errMsg.Code)
}

func (s *UserTestSuite) TestReadUser() {
	jsonBody, err := json.Marshal(s.user)
	s.NoError(err)
	r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Get user by email.
	r = httptest.NewRequest("GET", "/api/v1/users/by-email/test@yandex.ru", nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	var user UserMessage
	_ = json.NewDecoder(w.Body).Decode(&user)

	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(s.user.Email, user.Email)
	s.Equal(s.user.FirstName, user.FirstName)
	s.Equal(s.user.LastName, user.LastName)
}

func (s *UserTestSuite) TestUpdateUser() {
	jsonBody, err := json.Marshal(s.user)
	s.NoError(err)
	r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Modify exist user.
	s.user.FirstName = "Test Username"
	jsonBody, err = json.Marshal(s.user)
	s.NoError(err)
	r = httptest.NewRequest("PUT", "/api/v1/users/1", bytes.NewReader(jsonBody))

	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	s.NoError(err)

	var user UserMessage
	_ = json.NewDecoder(w.Body).Decode(&user)
	s.Equal("Test Username", user.FirstName)
}

func (s *UserTestSuite) TestDeleteUser() {
	jsonBody, err := json.Marshal(s.user)
	s.NoError(err)
	r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Delete user.
	r = httptest.NewRequest("DELETE", "/api/v1/users/1", nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	var user UserMessage
	_ = json.NewDecoder(w.Body).Decode(&user)
	s.NoError(err)

	s.Equal(http.StatusOK, w.Result().StatusCode)

	// Check that user not exist.
	r = httptest.NewRequest("GET", "/api/v1/users/by-email/test@yandex.ru", nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	var errMsg ErrorMessage
	err = json.NewDecoder(w.Body).Decode(&errMsg)
	s.NoError(err)

	s.Equal(http.StatusInternalServerError, w.Result().StatusCode)
	s.Equal("user doesn't not exist", errMsg.Message)
}

type EventTestSuite struct {
	suite.Suite
	mux   *runtime.ServeMux
	event structs.Event
}

func TestEventSuite(t *testing.T) {
	// This is what actually runs our suite
	suite.Run(t, new(EventTestSuite))
}

func (s *EventTestSuite) SetupTest() {
	cfg := config.NewConfig("testdata/config.yaml")

	// Init in-memory database.
	database := cfg.DatabaseBuilder()
	err := database.Connect(cfg.GetDatabasePath())
	s.NoError(err)

	// Default defined event.
	s.event = structs.Event{
		UserID:   1,
		Title:    "Base title",
		Content:  "Random content",
		DateFrom: time.Date(2020, 1, 2, 12, 4, 37, 0, time.UTC),
		DateTo:   time.Date(2020, 1, 3, 9, 15, 0, 0, time.UTC),
		Notified: false,
	}

	// Init server.
	ctx, cancel := context.WithCancel(context.Background())
	s.mux = runtime.NewServeMux()
	srv := &Server{
		ctx:    ctx,
		cancel: cancel,

		db: database,
	}

	// Register server.
	RegisterCalendarHandlerServer(context.TODO(), s.mux, srv)
}

func (s *EventTestSuite) TestEmptyEvent() {
	r := httptest.NewRequest("GET", "/api/v1/events", nil)
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	s.Equal(http.StatusMethodNotAllowed, w.Result().StatusCode)
}

func (s *EventTestSuite) TestCreateEvent() {
	// Create user.
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	var eventMsg EventMessage
	err = json.NewDecoder(w.Body).Decode(&eventMsg)
	s.NoError(err)

	e, err := structs.ConvertEventMessageToEvent(eventMsg)
	s.NoError(err)

	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(s.event.Title, eventMsg.Title)
	s.Equal(s.event.Content, eventMsg.Content)
	s.Equal(s.event.DateFrom, e.DateFrom)
	s.Equal("2020-01-02T12:04:37Z", e.DateFrom.Format(time.RFC3339))
}

func (s *EventTestSuite) TestReadEvents() {
	// Create first event.
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Create second event.
	s.event.DateFrom = s.event.DateFrom.Add(2 * time.Hour)
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)

	r = httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Read events.
	r = httptest.NewRequest("GET", "/api/v1/events/by-user/1", nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	var manyEvents ManyEventsMessage
	err = json.NewDecoder(w.Body).Decode(&manyEvents)
	s.NoError(err)

	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(2, len(manyEvents.Events))

	e, err := structs.ConvertEventMessageToEvent(manyEvents.Events[0])
	s.NoError(err)
	s.Equal("2020-01-02T12:04:37Z", e.DateFrom.Format(time.RFC3339))

	e, err = structs.ConvertEventMessageToEvent(manyEvents.Events[1])
	s.NoError(err)
	s.Equal("2020-01-02T14:04:37Z", e.DateFrom.Format(time.RFC3339))
}

func (s *EventTestSuite) TestUpdateEvents() {
	// Create event.
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	var eventMsg EventMessage
	err = json.NewDecoder(w.Body).Decode(&eventMsg)
	s.NoError(err)
	e, err := structs.ConvertEventMessageToEvent(eventMsg)
	s.NoError(err)
	s.Equal("2020-01-02T12:04:37Z", e.DateFrom.Format(time.RFC3339))

	// Update event.
	s.event.DateFrom = s.event.DateFrom.Add(2 * time.Hour)
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)

	r = httptest.NewRequest("PUT", "/api/v1/events/1", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	err = json.NewDecoder(w.Body).Decode(&eventMsg)
	s.NoError(err)
	e, err = structs.ConvertEventMessageToEvent(eventMsg)
	s.NoError(err)
	s.Equal("2020-01-02T14:04:37Z", e.DateFrom.Format(time.RFC3339))
}

func (s *EventTestSuite) TestDeleteEvent() {
	// Create event.
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Delete event.
	r = httptest.NewRequest("DELETE", "/api/v1/events/1", nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	var eventMsg EventMessage
	err = json.NewDecoder(w.Body).Decode(&eventMsg)
	s.NoError(err)
	e, err := structs.ConvertEventMessageToEvent(eventMsg)
	s.NoError(err)
	s.Equal(int64(1), e.ID)

	// Check no events
	r = httptest.NewRequest("GET", "/api/v1/events/by-user/1", nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	var manyEvents ManyEventsMessage
	err = json.NewDecoder(w.Body).Decode(&manyEvents)
	s.NoError(err)

	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(0, len(manyEvents.Events))
}

func (s *EventTestSuite) TestRangeEvents() {
	// First event.
	s.event.DateFrom = time.Date(2020, 10, 1, 12, 0, 0, 0, time.UTC)
	s.event.DateTo = s.event.DateFrom.Add(+2 * time.Hour)
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Second event
	s.event.DateFrom = time.Date(2020, 9, 30, 12, 0, 0, 0, time.UTC)
	s.event.DateTo = s.event.DateFrom.Add(+2 * time.Hour)
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)
	r = httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Third event
	s.event.DateFrom = time.Date(2020, 9, 30, 16, 0, 0, 0, time.UTC)
	s.event.DateTo = s.event.DateFrom.Add(+2 * time.Hour)
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)
	r = httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Fourth event
	s.event.DateFrom = time.Date(2020, 9, 15, 16, 0, 0, 0, time.UTC)
	s.event.DateTo = s.event.DateFrom.Add(+2 * time.Hour)
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)
	r = httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Result - 4 events summary:
	// - 1 event in October
	// - 3 event in September
	// - 2 event in daily 30.09
	// - 1 event in day 15.09

	// Check summery events
	r = httptest.NewRequest("GET", "/api/v1/events/by-user/1", nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)
	manyEvents := new(ManyEventsMessage)
	err = json.NewDecoder(w.Body).Decode(manyEvents)
	s.NoError(err)
	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(4, len(manyEvents.Events))

	// Check October events
	eventsURL := "/api/v1/events/by-date/1/date/" + url.QueryEscape("2020-10-20T12:04:05+03:00") + "/period/MONTHLY"
	r = httptest.NewRequest("GET", eventsURL, nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)
	manyEvents = new(ManyEventsMessage)
	err = json.NewDecoder(w.Body).Decode(&manyEvents)
	s.NoError(err)
	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(1, len(manyEvents.Events))

	// Check September events
	eventsURL = "/api/v1/events/by-date/1/date/" + url.QueryEscape("2020-09-10T05:00:00+03:00") + "/period/MONTHLY"
	r = httptest.NewRequest("GET", eventsURL, nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)
	manyEvents = new(ManyEventsMessage)
	err = json.NewDecoder(w.Body).Decode(&manyEvents)
	s.NoError(err)
	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(3, len(manyEvents.Events))

	// Check week 28.09 - 04.10 events
	eventsURL = "/api/v1/events/by-date/1/date/" + url.QueryEscape("2020-09-30T08:00:01+03:00") + "/period/WEEKLY"
	r = httptest.NewRequest("GET", eventsURL, nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)
	manyEvents = new(ManyEventsMessage)
	err = json.NewDecoder(w.Body).Decode(&manyEvents)
	s.NoError(err)
	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(3, len(manyEvents.Events))

	// Check daily event 2020.09.15
	eventsURL = "/api/v1/events/by-date/1/date/" + url.QueryEscape("2020-09-15T08:00:01+03:00") + "/period/DAILY"
	r = httptest.NewRequest("GET", eventsURL, nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)
	manyEvents = new(ManyEventsMessage)
	err = json.NewDecoder(w.Body).Decode(&manyEvents)
	s.NoError(err)
	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(1, len(manyEvents.Events))
}

func (s *EventTestSuite) TestReadUninformedEvents() {
	s.event.DateFrom = time.Now().Add(+20 * time.Hour * 20)
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Check no events
	r = httptest.NewRequest("GET", "/api/v1/events/uninformed", nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	var manyEvents ManyEventsMessage
	err = json.NewDecoder(w.Body).Decode(&manyEvents)
	s.NoError(err)

	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(0, len(manyEvents.Events))

	// Create current event
	s.event.DateFrom = time.Now()
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)

	r = httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	// Check one events
	r = httptest.NewRequest("GET", "/api/v1/events/uninformed", nil)
	w = httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)

	err = json.NewDecoder(w.Body).Decode(&manyEvents)
	s.NoError(err)

	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal(1, len(manyEvents.Events))
}
