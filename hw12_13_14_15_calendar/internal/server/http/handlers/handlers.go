package internalhttphandlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
)

// Possible statuses.
const (
	StatusOK   = "successful"
	StatusFail = "failed"
)

type ErrorBody struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func generateError(w io.Writer, err error) error {
	return json.NewEncoder(w).Encode(ErrorBody{
		Status: StatusFail,
		Error:  fmt.Sprintf("%s", err),
	})
}

type ServeHandler struct {
	db  storage.ClientInterface
	log *logger.Logger
}

func NewServeHandler(log *logger.Logger, database storage.ClientInterface) *ServeHandler {
	srvHander := &ServeHandler{
		db:  database,
		log: log,
	}

	return srvHander
}

func (sh *ServeHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Alive"))
}
