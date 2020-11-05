package logger

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/stretchr/testify/require"
)

type LoggerParams struct {
	Output   []string `yaml:"output"`
	Error    []string `yaml:"error"`
	Level    string   `yaml:"level"`
	Encoding string   `yaml:"encoding"`
}

func TestLogger(t *testing.T) {
	t.Run("create logger", func(t *testing.T) {
		defer os.Remove("stdout_file")
		_ = os.Remove("stdout_file")

		log := NewLogger(&config.Configuration{
			Logger: config.LoggerParams{
				Output: []string{"stdout_file"},
				Level:  "info",
			},
		})
		require.NotEmpty(t, log)

		// Write data to logfile.
		log.Error("check file for error")
		log.Info("information message")

		// This message will be ignored.
		log.Debug("Debug message")

		// Close logger object.
		err := log.Close()
		require.NoError(t, err)

		// Read changed file.
		logData, err := ioutil.ReadFile("stdout_file")
		require.NoError(t, err)
		stringData := string(logData)
		require.Equal(t, strings.Count(stringData, "Debug message"), 0)
		require.NotEqual(t, strings.Count(stringData, "information message"), 0)
		require.NotEqual(t, strings.Count(stringData, "check file for error"), 0)
	})
}
