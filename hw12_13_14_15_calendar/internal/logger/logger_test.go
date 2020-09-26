package logger

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	config "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
)

func TestLogger(t *testing.T) {
	t.Run("create logger", func(t *testing.T) {
		defer os.Remove("stdout_file")
		_ = os.Remove("stdout_file")

		log := NewLogger(&config.Configuration{
			Logger: config.LoggerParams{
				Output:   []string{"stdout_file"},
				Level:    "info",
				Encoding: "json",
			}})
		require.NotEmpty(t, log)

		log.Error("check file for error")
		log.Info("information message")

		// this message must be ignored
		log.Debug("Debug message")

		require.Nil(t, log.Close())

		logData, err := ioutil.ReadFile("stdout_file")
		require.Nil(t, err)

		stringData := string(logData)
		require.Equal(t, strings.Count(stringData, "Debug message"), 0)
		require.NotEqual(t, strings.Count(stringData, "information message"), 0)
		require.NotEqual(t, strings.Count(stringData, "check file for error"), 0)
	})
}
