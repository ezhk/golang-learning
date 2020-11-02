package godogs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cucumber/godog"
	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/structs"
)

var (
	Server       string
	ResponseBody string

	EventID int
	Event   structs.EventMessage
	Events  structs.ManyEventsMessage
)

func checkEventNotificationState() error {
	err := json.Unmarshal([]byte(ResponseBody), &Events)
	if err != nil {
		return err
	}

	return nil
}

func createEventWithUserIdTitleDateFromDateTo(arg1 int, arg2, arg3, arg4 string) error {
	client := &http.Client{}
	event := struct {
		UserID   int    `json:"UserID"`
		Title    string `json:"Title"`
		DateFrom string `json:"DateFrom"`
		DateTo   string `json:"DateTo"`
		Notified bool   `json:"Notified"`
	}{
		UserID:   arg1,
		Title:    arg2,
		DateFrom: arg3,
		DateTo:   arg4,
		Notified: false,
	}

	json, err := json.Marshal(event)
	if err != nil {
		return err
	}

	URL := fmt.Sprintf("%s/api/v1/events", Server)
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	ResponseBody = buf.String()

	return nil
}

func createEventWithUserIdTitleStartsInDays(arg1 int, arg2 string, arg3 int) error {
	client := &http.Client{}
	startTime := time.Now().Add(time.Hour * time.Duration(24*arg3))
	endTime := startTime.Add(time.Hour * 1)
	event := struct {
		UserID   int    `json:"UserID"`
		Title    string `json:"Title"`
		DateFrom string `json:"DateFrom"`
		DateTo   string `json:"DateTo"`
	}{
		UserID:   arg1,
		Title:    arg2,
		DateFrom: startTime.Format(time.RFC3339),
		DateTo:   endTime.Format(time.RFC3339),
	}

	json, err := json.Marshal(event)
	if err != nil {
		return err
	}

	URL := fmt.Sprintf("%s/api/v1/events", Server)
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	ResponseBody = buf.String()

	return nil
}

func createUserWithEmail(arg1 string) error {
	client := &http.Client{}
	user := struct {
		Email string `json:"Email"`
	}{Email: arg1}

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	URL := fmt.Sprintf("%s/api/v1/users", Server)
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	ResponseBody = buf.String()

	return nil
}

func deleteEventsForUserId(arg1 int) error {
	client := &http.Client{}

	// Fetch all user events.
	err := getEventsByUserId(arg1)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(ResponseBody), &Events)
	if err != nil {
		return err
	}

	for _, e := range Events.Events {
		URL := fmt.Sprintf("%s/api/v1/events/%s", Server, e.ID)
		req, err := http.NewRequest(http.MethodDelete, URL, nil)
		if err != nil {
			return err
		}

		_, err = client.Do(req)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteEventByGlobalID() error {
	client := &http.Client{}
	URL := fmt.Sprintf("%s/api/v1/events/%d", Server, EventID)
	req, err := http.NewRequest(http.MethodDelete, URL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	ResponseBody = buf.String()

	return nil
}

func deleteEventByID(arg1 int) error {
	EventID = arg1

	return deleteEventByGlobalID()
}

func existEventWithGlobalID() error {
	err := json.Unmarshal([]byte(ResponseBody), &Events)
	if err != nil {
		return err
	}

	for _, e := range Events.Events {
		if e.ID == strconv.Itoa(EventID) {
			return nil
		}
	}

	return errors.New("event doesn't exist")
}

func getEventsByUserId(arg1 int) error {
	URL := fmt.Sprintf("%s/api/v1/events/by-user/%d", Server, arg1)
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	ResponseBody = buf.String()

	return nil
}

func getEventsByUserIdAfterSeconds(arg1, arg2 int) error {
	time.Sleep(time.Duration(arg2) * time.Second)
	URL := fmt.Sprintf("%s/api/v1/events/by-user/%d", Server, arg1)
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	ResponseBody = buf.String()

	return nil
}

func getEventsByUserIdDateAndPeriod(arg1 int, arg2, arg3 string) error {
	URL := fmt.Sprintf("%s/api/v1/events/by-date/%d/date/%s/period/%s", Server, arg1, arg2, arg3)
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	ResponseBody = buf.String()

	return nil
}

func getUserByEmail(arg1 string) error {
	URL := fmt.Sprintf("%s/api/v1/users/by-email/%s", Server, arg1)
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	ResponseBody = buf.String()

	return nil
}

func nonExistEventWithID(arg1 int) error {
	Event = structs.EventMessage{
		ID: strconv.Itoa(arg1),
	}

	return nil
}

func receiveErrorDeleteStatus() error {
	var errorMessage structs.ErrorMessage
	err := json.Unmarshal([]byte(ResponseBody), &errorMessage)
	if err != nil {
		return err
	}

	return nil
}

func receiveSuccessfulDeleteStatus() error {
	err := json.Unmarshal([]byte(ResponseBody), &Event)
	if err != nil {
		return err
	}

	return nil
}

func receiveSuccessfulListEventsWithEvents(arg1 int) error {
	err := json.Unmarshal([]byte(ResponseBody), &Events)
	if err != nil {
		return err
	}

	if len(Events.Events) != arg1 {
		return errors.New("incorrent events count")
	}

	return nil
}

func receiveSuccessfulSaveStatus() error {
	err := json.Unmarshal([]byte(ResponseBody), &Event)
	if err != nil {
		return err
	}

	return nil
}

func receiveSuccessfulUpdateStatus() error {
	err := json.Unmarshal([]byte(ResponseBody), &Event)
	if err != nil {
		return err
	}

	return nil
}

func receivedEmptyUserAnswer() error {
	var errorMessage structs.ErrorMessage
	err := json.Unmarshal([]byte(ResponseBody), &errorMessage)
	if err != nil {
		return err
	}

	return nil
}

func receivedFalseStatus() error {
	for _, e := range Events.Events {
		if e.Notified {
			return errors.New("exist notified messages")
		}
	}

	return nil
}

func receivedNonEmptyUserAnswer() error {
	var user structs.UserMessage
	err := json.Unmarshal([]byte(ResponseBody), &user)
	if err != nil {
		return err
	}

	return nil
}

func receivedTrueStatus() error {
	for _, e := range Events.Events {
		if e.Notified {
			return nil
		}
	}

	return errors.New("no notified messages")
}

func saveEventIDAsGlobalID() error {
	err := json.Unmarshal([]byte(ResponseBody), &Events)
	if err != nil {
		return err
	}

	// Get first event as based.
	EventID, err = strconv.Atoi(Events.Events[0].ID)
	if err != nil {
		return err
	}

	return nil
}

func updateEventWithGlobalIDAndTitle(arg1 string) error {
	Event.Title = arg1

	client := &http.Client{}
	json, err := json.Marshal(Event)
	if err != nil {
		return err
	}

	URL := fmt.Sprintf("%s/api/v1/events/%s", Server, Event.ID)
	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	ResponseBody = buf.String()

	return nil
}

func FeatureContext(s *godog.Suite) {
	s.BeforeScenario(func(*godog.Scenario) {
		Server = os.Getenv("REST_SERVER")

		if len(Server) < 1 {
			Server = "http://localhost:8080"
		}
	})

	// func InitializeScenario(s *godog.ScenarioContext) {
	s.Step(`^check event notification state$`, checkEventNotificationState)
	s.Step(`^create event with user id (\d+), title "([^"]*)", date from "([^"]*)", date to "([^"]*)"$`, createEventWithUserIdTitleDateFromDateTo)
	s.Step(`^create event with user id (\d+), title "([^"]*)", starts in (\d+) days$`, createEventWithUserIdTitleStartsInDays)
	s.Step(`^create user with email "([^"]*)"$`, createUserWithEmail)
	s.Step(`^delete events for user id (\d+)$`, deleteEventsForUserId)
	s.Step(`^delete event by global ID$`, deleteEventByGlobalID)
	s.Step(`^delete event by ID (\d+)$`, deleteEventByID)
	s.Step(`^exist event with global ID$`, existEventWithGlobalID)
	s.Step(`^get events by user id (\d+)$`, getEventsByUserId)
	s.Step(`^get events by user id (\d+) after (\d+) seconds$`, getEventsByUserIdAfterSeconds)
	s.Step(`^get events by user id (\d+), date "([^"]*)" and period "([^"]*)"$`, getEventsByUserIdDateAndPeriod)
	s.Step(`^get user by email "([^"]*)"$`, getUserByEmail)
	s.Step(`^non exist event with ID (\d+)$`, nonExistEventWithID)
	s.Step(`^receive error delete status$`, receiveErrorDeleteStatus)
	s.Step(`^receive successful delete status$`, receiveSuccessfulDeleteStatus)
	s.Step(`^receive successful list events with (\d+) events$`, receiveSuccessfulListEventsWithEvents)
	s.Step(`^receive successful save status$`, receiveSuccessfulSaveStatus)
	s.Step(`^receive successful update status$`, receiveSuccessfulUpdateStatus)
	s.Step(`^received empty user answer$`, receivedEmptyUserAnswer)
	s.Step(`^received false status$`, receivedFalseStatus)
	s.Step(`^received non empty user answer$`, receivedNonEmptyUserAnswer)
	s.Step(`^received true status$`, receivedTrueStatus)
	s.Step(`^save event ID as global ID$`, saveEventIDAsGlobalID)
	s.Step(`^update event with global ID and title "([^"]*)"$`, updateEventWithGlobalIDAndTitle)
}

// var opts = godog.Options{
// 	Output: colors.Colored(os.Stdout),
// 	Format: "progress",
// }

// func init() {
// 	Server = os.Getenv("REST_SERVER")

// 	if len(Server) < 1 {
// 		Server = "http://localhost:8080"
// 	}
// }

// func TestMain(m *testing.M) {
// 	// flag.Parse()
// 	// opts.Paths = flag.Args()

// 	status := godog.TestSuite{
// 		Name: "godogs",
// 		// TestSuiteInitializer: InitializeTestSuite,
// 		ScenarioInitializer: InitializeScenario,
// 		Options:             &opts,
// 	}.Run()

// 	fmt.Println(status)
// 	os.Exit(status)
// }
