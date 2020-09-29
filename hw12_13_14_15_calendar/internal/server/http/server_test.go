package internalhttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	internalhttphandlers "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/server/http/handlers"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
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
