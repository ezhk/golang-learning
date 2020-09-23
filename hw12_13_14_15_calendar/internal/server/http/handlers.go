package internalhttp

import "net/http"

type ServeHandler struct{}

func (sh *ServeHandler) Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("world"))
}
