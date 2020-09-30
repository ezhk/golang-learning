package internalhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	internalhttphandlers "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/server/http/handlers"
	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	user    storage.User
	handler *internalhttphandlers.ServeHandler
}

type EventTestSuite struct {
	suite.Suite
	event   storage.Event
	handler *internalhttphandlers.ServeHandler
}

func TestHealthCheck(t *testing.T) {
	t.Run("health-check handler", func(t *testing.T) {
		handler := &internalhttphandlers.ServeHandler{}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()
		handler.HealthCheck(w, r)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.Nil(t, err)
		require.Equal(t, []byte("Alive"), body)
	})
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

	s.handler = internalhttphandlers.NewServeHandler(nil, database)
	s.user = storage.User{
		Email:     "test@yandex.ru",
		FirstName: "Test",
		LastName:  "Yandex mail",
	}
}

func (s *UserTestSuite) TestEmptyUser() {
	r := httptest.NewRequest("GET", "/users", nil)
	// Set URL param is not simple URL path.
	r = mux.SetURLVars(r, map[string]string{"email": "test@yandex.ru"})
	w := httptest.NewRecorder()
	s.handler.ReadUser(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var userError internalhttphandlers.ErrorBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&userError)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusFail, userError.Status)
	s.Equal("user doesn't not exist", userError.Error)
}

func (s *UserTestSuite) TestCreateUser() {
	// Create user.
	jsonBody, err := json.Marshal(s.user)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.handler.CreateUser(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var userBody internalhttphandlers.UserBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&userBody)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusOK, userBody.Status)
	s.Equal(s.user.Email, userBody.User.Email)
	s.Equal(s.user.FirstName, userBody.User.FirstName)
	s.Equal(s.user.LastName, userBody.User.LastName)

	// Create again might cause error.
	r = httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.handler.CreateUser(w, r)

	body, err = ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var userError internalhttphandlers.ErrorBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&userError)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusFail, userError.Status)
	s.Equal("user already exists", userError.Error)
}

func (s *UserTestSuite) TestReadUser() {
	jsonBody, err := json.Marshal(s.user)
	s.NoError(err)
	r := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.handler.CreateUser(w, r)

	// Get user by email.
	r = httptest.NewRequest("GET", "/users/test@yandex.ru", nil)
	// Set URL param is not simple URL path.
	r = mux.SetURLVars(r, map[string]string{"email": "test@yandex.ru"})
	w = httptest.NewRecorder()
	s.handler.ReadUser(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var userBody internalhttphandlers.UserBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&userBody)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusOK, userBody.Status)
	s.Equal(s.user.Email, userBody.User.Email)
	s.Equal(s.user.FirstName, userBody.User.FirstName)
	s.Equal(s.user.LastName, userBody.User.LastName)
}

func (s *UserTestSuite) TestUpdateUser() {
	jsonBody, err := json.Marshal(s.user)
	s.NoError(err)
	r := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.handler.CreateUser(w, r)

	// Modify exist user.
	s.user.FirstName = "Test Username"
	jsonBody, err = json.Marshal(s.user)
	s.NoError(err)
	r = httptest.NewRequest("PUT", "/users/1", bytes.NewReader(jsonBody))
	r = mux.SetURLVars(r, map[string]string{"id": "1"})

	w = httptest.NewRecorder()
	s.handler.UpdateUser(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	userBody := internalhttphandlers.UserBody{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&userBody)
	s.NoError(err)
	s.Equal("Test Username", userBody.User.FirstName)
}

func (s *UserTestSuite) TestDeleteUser() {
	jsonBody, err := json.Marshal(s.user)
	s.NoError(err)
	r := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.handler.CreateUser(w, r)

	// Delete user.
	r = httptest.NewRequest("DELETE", "/users/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	w = httptest.NewRecorder()
	s.handler.DeleteUser(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var userBody internalhttphandlers.UserBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&userBody)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusOK, userBody.Status)
	s.Equal(int64(1), userBody.User.ID)

	// Check that user not exist.
	r = httptest.NewRequest("GET", "/users", nil)
	r = mux.SetURLVars(r, map[string]string{"email": "test@yandex.ru"})
	w = httptest.NewRecorder()
	s.handler.ReadUser(w, r)

	body, err = ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var userError internalhttphandlers.ErrorBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&userError)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusFail, userError.Status)
	s.Equal("user doesn't not exist", userError.Error)
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

	s.handler = internalhttphandlers.NewServeHandler(nil, database)
	s.event = storage.Event{
		UserID:   1,
		Title:    "Base title",
		Content:  "Random content",
		DateFrom: time.Date(2020, 1, 2, 12, 4, 37, 0, time.UTC),
		DateTo:   time.Date(2020, 1, 3, 9, 15, 0, 0, time.UTC),
		Notified: false,
	}
}

func (s *EventTestSuite) TestEmptyEvent() {
	r := httptest.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()
	s.handler.ReadEvents(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var eventError internalhttphandlers.ErrorBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&eventError)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusFail, eventError.Status)
	s.Equal("empty events list", eventError.Error)
}

func (s *EventTestSuite) TestCreateEvent() {
	// Create user.
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var oneEventBody internalhttphandlers.OneEventBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&oneEventBody)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusOK, oneEventBody.Status)
	s.Equal(s.event.Title, oneEventBody.Event.Title)
	s.Equal(s.event.Content, oneEventBody.Event.Content)
	s.Equal(s.event.DateFrom, oneEventBody.Event.DateFrom)
	s.Equal("2020-01-02T12:04:37Z", oneEventBody.Event.DateFrom.Format(time.RFC3339))
}

func (s *EventTestSuite) TestReadEvents() {
	// Create first event.
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	// Create second event.
	s.event.DateFrom = s.event.DateFrom.Add(2 * time.Hour)
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)

	r = httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	// Read events.
	r = httptest.NewRequest("GET", "/events?userid=1", nil)
	w = httptest.NewRecorder()
	s.handler.ReadEvents(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var manyEventsBody internalhttphandlers.ManyEventsBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&manyEventsBody)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusOK, manyEventsBody.Status)
	s.Equal(2, len(manyEventsBody.Events))
	s.Equal("2020-01-02T12:04:37Z", manyEventsBody.Events[0].DateFrom.Format(time.RFC3339))
	s.Equal("2020-01-02T14:04:37Z", manyEventsBody.Events[1].DateFrom.Format(time.RFC3339))
}

func (s *EventTestSuite) TestUpdateEvents() {
	// Create event.
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var oneEventBody internalhttphandlers.OneEventBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&oneEventBody)
	s.NoError(err)
	s.Equal("2020-01-02T12:04:37Z", oneEventBody.Event.DateFrom.Format(time.RFC3339))

	// Update event.
	s.event.DateFrom = s.event.DateFrom.Add(2 * time.Hour)
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)

	r = httptest.NewRequest("PUT", "/events/1", bytes.NewReader(jsonBody))
	// Set URL param is not simple URL path.
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	w = httptest.NewRecorder()
	s.handler.UpdateEvent(w, r)

	body, err = ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	err = json.NewDecoder(bytes.NewReader(body)).Decode(&oneEventBody)
	s.NoError(err)
	s.Equal("2020-01-02T14:04:37Z", oneEventBody.Event.DateFrom.Format(time.RFC3339))
}

func (s *EventTestSuite) TestDeleteEvent() {
	// Create event.
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	// Delete event.
	r = httptest.NewRequest("DELETE", "/events/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	w = httptest.NewRecorder()
	s.handler.DeleteEvent(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var oneEventBody internalhttphandlers.OneEventBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&oneEventBody)
	s.NoError(err)
	s.Equal(storage.Event{ID: 1}, oneEventBody.Event)

	// Check no events
	r = httptest.NewRequest("GET", "/events?userid=1", nil)
	w = httptest.NewRecorder()
	s.handler.ReadEvents(w, r)

	body, err = ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var manyEventsBody internalhttphandlers.ManyEventsBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&manyEventsBody)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusOK, manyEventsBody.Status)
	s.Equal(0, len(manyEventsBody.Events))
}

func (s *EventTestSuite) TestRangeEvents() {
	// First event.
	s.event.DateFrom = time.Date(2020, 10, 1, 12, 0, 0, 0, time.UTC)
	s.event.DateTo = s.event.DateFrom.Add(+2 * time.Hour)
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	// Second event
	s.event.DateFrom = time.Date(2020, 9, 30, 12, 0, 0, 0, time.UTC)
	s.event.DateTo = s.event.DateFrom.Add(+2 * time.Hour)
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)
	r = httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	// Third event
	s.event.DateFrom = time.Date(2020, 9, 30, 16, 0, 0, 0, time.UTC)
	s.event.DateTo = s.event.DateFrom.Add(+2 * time.Hour)
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)
	r = httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	// Fourth event
	s.event.DateFrom = time.Date(2020, 9, 15, 16, 0, 0, 0, time.UTC)
	s.event.DateTo = s.event.DateFrom.Add(+2 * time.Hour)
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)
	r = httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	// Result - 4 events summary:
	// - 1 event in October
	// - 3 event in September
	// - 2 event in daily 30.09
	// - 1 event in day 15.09

	// Check summery events
	r = httptest.NewRequest("GET", "/events?userid=1", nil)
	w = httptest.NewRecorder()
	s.handler.ReadEvents(w, r)
	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)
	var manyEventsBody internalhttphandlers.ManyEventsBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&manyEventsBody)
	s.NoError(err)
	s.Equal(internalhttphandlers.StatusOK, manyEventsBody.Status)
	s.Equal(4, len(manyEventsBody.Events))

	// Check October events
	eventsURL := "/events?userid=1&filter=monthly&date=" + url.QueryEscape("2020-10-20T12:04:05+03:00")
	r = httptest.NewRequest("GET", eventsURL, nil)
	w = httptest.NewRecorder()
	s.handler.ReadEvents(w, r)
	body, err = ioutil.ReadAll(w.Result().Body)
	s.NoError(err)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&manyEventsBody)
	s.NoError(err)
	s.Equal(internalhttphandlers.StatusOK, manyEventsBody.Status)
	s.Equal(1, len(manyEventsBody.Events))

	// Check September events
	eventsURL = "/events?userid=1&filter=monthly&date=" + url.QueryEscape("2020-09-10T05:00:00+03:00")
	r = httptest.NewRequest("GET", eventsURL, nil)
	w = httptest.NewRecorder()
	s.handler.ReadEvents(w, r)
	body, err = ioutil.ReadAll(w.Result().Body)
	s.NoError(err)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&manyEventsBody)
	s.NoError(err)
	s.Equal(internalhttphandlers.StatusOK, manyEventsBody.Status)
	s.Equal(3, len(manyEventsBody.Events))

	// Check week 28.09 - 04.10 events
	eventsURL = "/events?userid=1&filter=daily&date=" + url.QueryEscape("2020-09-30T08:00:01+03:00")
	r = httptest.NewRequest("GET", eventsURL, nil)
	w = httptest.NewRecorder()
	s.handler.ReadEvents(w, r)
	body, err = ioutil.ReadAll(w.Result().Body)
	s.NoError(err)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&manyEventsBody)
	s.NoError(err)
	s.Equal(internalhttphandlers.StatusOK, manyEventsBody.Status)
	s.Equal(2, len(manyEventsBody.Events))

	// Check daily event 2020.09.15
	eventsURL = "/events?userid=1&filter=daily&date=" + url.QueryEscape("2020-09-15T08:00:01+03:00")
	r = httptest.NewRequest("GET", eventsURL, nil)
	w = httptest.NewRecorder()
	s.handler.ReadEvents(w, r)
	body, err = ioutil.ReadAll(w.Result().Body)
	s.NoError(err)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&manyEventsBody)
	s.NoError(err)
	s.Equal(internalhttphandlers.StatusOK, manyEventsBody.Status)
	s.Equal(1, len(manyEventsBody.Events))
}

func (s *EventTestSuite) TestReadUninformedEvents() {
	s.event.DateFrom = time.Now().Add(+20 * time.Hour * 20)
	jsonBody, err := json.Marshal(s.event)
	s.NoError(err)

	r := httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	// Check no events
	r = httptest.NewRequest("GET", "/events?filter=uninformed", nil)
	w = httptest.NewRecorder()
	s.handler.ReadEvents(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	var manyEventsBody internalhttphandlers.ManyEventsBody
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&manyEventsBody)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusOK, manyEventsBody.Status)
	s.Equal(0, len(manyEventsBody.Events))

	// Create current event
	s.event.DateFrom = time.Now()
	jsonBody, err = json.Marshal(s.event)
	s.NoError(err)

	r = httptest.NewRequest("POST", "/events", bytes.NewReader(jsonBody))
	w = httptest.NewRecorder()
	s.handler.CreateEvent(w, r)

	// Check one events
	r = httptest.NewRequest("GET", "/events?filter=uninformed", nil)
	w = httptest.NewRecorder()
	s.handler.ReadEvents(w, r)

	body, err = ioutil.ReadAll(w.Result().Body)
	s.NoError(err)

	err = json.NewDecoder(bytes.NewReader(body)).Decode(&manyEventsBody)
	s.NoError(err)

	s.Equal(internalhttphandlers.StatusOK, manyEventsBody.Status)
	s.Equal(1, len(manyEventsBody.Events))
}
