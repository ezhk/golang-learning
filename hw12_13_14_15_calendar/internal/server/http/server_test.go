package internalhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	internalhttphandlers "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/server/http/handlers"
	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

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

func TestUser(t *testing.T) {
	t.Run("empty user error", func(t *testing.T) {
		handler := internalhttphandlers.NewServeHandler(&config.Configuration{DB: config.DatabaseParams{Type: "in-memory"}}, nil)

		r := httptest.NewRequest("GET", "/users", nil)
		// Set URL param is not simple URL path.
		r = mux.SetURLVars(r, map[string]string{"email": "test@yandex.ru"})
		w := httptest.NewRecorder()
		handler.ReadUser(w, r)

		body, err := ioutil.ReadAll(w.Result().Body)
		require.Nil(t, err)

		var userError internalhttphandlers.ErrorBody
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&userError)
		require.NoError(t, err)

		require.Equal(t, internalhttphandlers.StatusFail, userError.Status)
		require.Equal(t, "user doesn't not exist", userError.Error)
	})

	t.Run("create user", func(t *testing.T) {
		handler := internalhttphandlers.NewServeHandler(&config.Configuration{DB: config.DatabaseParams{Type: "in-memory"}}, nil)

		// Create user.
		user := storage.User{
			Email:     "test@yandex.ru",
			FirstName: "Test",
			LastName:  "Yandex mail",
		}

		jsonBody, err := json.Marshal(user)
		require.NoError(t, err)

		r := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()
		handler.CreateUser(w, r)

		body, err := ioutil.ReadAll(w.Result().Body)
		require.NoError(t, err)

		var userBody internalhttphandlers.UserBody
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&userBody)
		require.NoError(t, err)

		require.Equal(t, internalhttphandlers.StatusOK, userBody.Status)
		require.Equal(t, "test@yandex.ru", userBody.User.Email)
		require.Equal(t, "Test", userBody.User.FirstName)
		require.Equal(t, "Yandex mail", userBody.User.LastName)

		// Create again might cause error.
		r = httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
		w = httptest.NewRecorder()
		handler.CreateUser(w, r)

		body, err = ioutil.ReadAll(w.Result().Body)
		require.NoError(t, err)

		var userError internalhttphandlers.ErrorBody
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&userError)
		require.NoError(t, err)

		require.Equal(t, internalhttphandlers.StatusFail, userError.Status)
		require.Equal(t, "user already exists", userError.Error)
	})

	t.Run("read user", func(t *testing.T) {
		handler := internalhttphandlers.NewServeHandler(&config.Configuration{DB: config.DatabaseParams{Type: "in-memory"}}, nil)

		// Create user.
		user := storage.User{
			Email:     "test@yandex.ru",
			FirstName: "Test",
			LastName:  "Yandex mail",
		}
		jsonBody, err := json.Marshal(user)
		require.NoError(t, err)
		r := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()
		handler.CreateUser(w, r)

		// Get user by email.
		r = httptest.NewRequest("GET", "/users/test@yandex.ru", nil)
		// Set URL param is not simple URL path.
		r = mux.SetURLVars(r, map[string]string{"email": "test@yandex.ru"})
		w = httptest.NewRecorder()
		handler.ReadUser(w, r)

		body, err := ioutil.ReadAll(w.Result().Body)
		require.Nil(t, err)

		var userBody internalhttphandlers.UserBody
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&userBody)
		require.NoError(t, err)

		require.Equal(t, internalhttphandlers.StatusOK, userBody.Status)
		require.Equal(t, "test@yandex.ru", userBody.User.Email)
		require.Equal(t, "Test", userBody.User.FirstName)
		require.Equal(t, "Yandex mail", userBody.User.LastName)
	})

	t.Run("update user", func(t *testing.T) {
		handler := internalhttphandlers.NewServeHandler(&config.Configuration{DB: config.DatabaseParams{Type: "in-memory"}}, nil)

		// Create user.
		user := storage.User{
			Email:     "test@yandex.ru",
			FirstName: "Test",
			LastName:  "Yandex mail",
		}
		jsonBody, err := json.Marshal(user)
		require.NoError(t, err)
		r := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()
		handler.CreateUser(w, r)

		// Modify exist user.
		user.FirstName = "Test Username"
		jsonBody, err = json.Marshal(user)
		require.NoError(t, err)
		r = httptest.NewRequest("PUT", "/users/1", bytes.NewReader(jsonBody))
		r = mux.SetURLVars(r, map[string]string{"id": "1"})

		w = httptest.NewRecorder()
		handler.UpdateUser(w, r)

		body, err := ioutil.ReadAll(w.Result().Body)
		require.NoError(t, err)

		userBody := internalhttphandlers.UserBody{}
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&userBody)
		require.NoError(t, err)
		require.Equal(t, "Test Username", userBody.User.FirstName)
	})

	t.Run("delete user", func(t *testing.T) {
		handler := internalhttphandlers.NewServeHandler(&config.Configuration{DB: config.DatabaseParams{Type: "in-memory"}}, nil)

		// Create user.
		user := storage.User{
			Email:     "test@yandex.ru",
			FirstName: "Test",
			LastName:  "Yandex mail",
		}
		jsonBody, err := json.Marshal(user)
		require.NoError(t, err)
		r := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()
		handler.CreateUser(w, r)

		// Delete user.
		r = httptest.NewRequest("DELETE", "/users/1", nil)
		r = mux.SetURLVars(r, map[string]string{"id": "1"})
		w = httptest.NewRecorder()
		handler.DeleteUser(w, r)

		body, err := ioutil.ReadAll(w.Result().Body)
		require.Nil(t, err)

		var userBody internalhttphandlers.UserBody
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&userBody)
		require.NoError(t, err)

		require.Equal(t, internalhttphandlers.StatusOK, userBody.Status)
		require.Equal(t, int64(1), userBody.User.ID)
	})
}
