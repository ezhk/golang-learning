package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Run("hello handler", func(t *testing.T) {
		handler := &ServeHandler{}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()
		handler.Hello(w, r)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.Nil(t, err)
		require.Equal(t, []byte("world"), body)
	})
}
