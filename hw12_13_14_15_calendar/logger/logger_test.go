package logger

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("create logger", func(t *testing.T) {
		defer os.Remove("stdout_file")
		_ = os.Remove("stdout_file")

		viper.Set("logger.output", []string{"stdout_file"})
		viper.Set("logger.level", "info")

		log := NewLogger()
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
