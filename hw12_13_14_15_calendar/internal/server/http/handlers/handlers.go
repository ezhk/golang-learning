package internalhttphandlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	config "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
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

func NewServeHandler(cfg *config.Configuration, log *logger.Logger) *ServeHandler {
	srvHander := &ServeHandler{
		db:  cfg.DatabaseBuilder(),
		log: log,
	}

	err := srvHander.db.Connect(cfg.DB.Path)
	if err != nil {
		log.Error("cannot conect to database: %v", err)
	}

	return srvHander
}

func (sh *ServeHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Alive"))
}
