package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("load config to struct", func(t *testing.T) {
		viper.SetConfigFile("testdata/config.yaml")
		err := viper.ReadInConfig()
		require.NoError(t, err)

		cfg := &Configuration{}
		err = viper.Unmarshal(cfg)
		require.NoError(t, err)
	})

	t.Run("zap config builder", func(t *testing.T) {
		cfg := NewConfig("testdata/config.yaml")
		zapConfig := cfg.ZapConfigBuilder()

		require.Equal(t, cfg.Logger.Encoding, zapConfig.Encoding)
		require.Equal(t, cfg.Logger.Output, zapConfig.OutputPaths)
		require.Equal(t, cfg.Logger.Error, zapConfig.ErrorOutputPaths)
	})

	t.Run("database config", func(t *testing.T) {
		randomTestString := "too long test DSN string with rand symbols: @&!ADFSAcxz#1"

		cfg := NewConfig("testdata/config.yaml")
		require.NotEqual(t, randomTestString, cfg.Database.Path)

		// Change DSN string.
		viper.Set("DSN", randomTestString)

		// Validate changed database path from DSN.
		modifedEnvCfg := NewConfig("testdata/config.yaml")
		databasePath := modifedEnvCfg.GetDatabasePath()
		require.Equal(t, randomTestString, databasePath)
		require.Equal(t, randomTestString, modifedEnvCfg.Database.Path)
	})

	t.Run("queue config", func(t *testing.T) {
		cfg := NewConfig("testdata/config.yaml")
		require.Equal(t, "localhost:6379", cfg.Queue.Path)
	})
}
