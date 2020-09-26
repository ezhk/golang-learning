package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("load config to struct", func(t *testing.T) {
		viper.SetConfigFile("../../configs/calendar.yaml")
		err := viper.ReadInConfig()
		require.NoError(t, err)

		var cfg Configuration
		err = viper.Unmarshal(&cfg)
		require.NoError(t, err)

		require.Equal(t, 42, cfg.Test.DoNotRemoveMe)
	})

	t.Run("new config", func(t *testing.T) {
		cfg := NewConfig("../../configs/calendar.yaml")
		require.Equal(t, 42, cfg.Test.DoNotRemoveMe)
	})

	t.Run("zap config builder", func(t *testing.T) {
		cfg := NewConfig("../../configs/calendar.yaml")
		zapConfig := cfg.ZapConfigBuilder()

		require.Equal(t, cfg.Logger.Encoding, zapConfig.Encoding)
		require.Equal(t, cfg.Logger.Output, zapConfig.OutputPaths)
		require.Equal(t, cfg.Logger.Error, zapConfig.ErrorOutputPaths)
	})

	t.Run("database config", func(t *testing.T) {
		randomTestString := "too long test DSN string with rand symbols: @&!ADFSAcxz#1"

		cfg := NewConfig("../../configs/calendar.yaml")
		require.NotEqual(t, randomTestString, cfg.DB.Path)

		// change DNS: works only for SQL type database.
		cfg.DB.Type = "sql"
		viper.Set("DSN", randomTestString)

		cfg.DatabaseConfigBuilder()
		require.Equal(t, randomTestString, cfg.DB.Path)
	})
}
