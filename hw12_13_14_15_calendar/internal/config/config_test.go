package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("load config to struct", func(t *testing.T) {
		viper.SetConfigFile("testdata/default.yaml")
		err := viper.ReadInConfig()
		require.NoError(t, err)

		var cfg Configuration
		err = viper.Unmarshal(&cfg)
		require.NoError(t, err)

		require.Equal(t, 42, cfg.Test.DoNotRemoveMe)
	})

	t.Run("new config", func(t *testing.T) {
		cfg := NewConfig("testdata/default.yaml")
		require.Equal(t, 42, cfg.Test.DoNotRemoveMe)
	})

	t.Run("zap config builder", func(t *testing.T) {
		cfg := NewConfig("testdata/default.yaml")
		zapConfig := cfg.ZapConfigBuilder()

		require.Equal(t, cfg.Logger.Encoding, zapConfig.Encoding)
		require.Equal(t, cfg.Logger.Output, zapConfig.OutputPaths)
		require.Equal(t, cfg.Logger.Error, zapConfig.ErrorOutputPaths)
	})

	t.Run("database config", func(t *testing.T) {
		randomTestString := "too long test DSN string with rand symbols: @&!ADFSAcxz#1"

		cfgInMemory := NewConfig("testdata/default.yaml")
		require.NotEqual(t, randomTestString, cfgInMemory.DB.Path)

		// change DNS: works only for SQL type database.
		viper.Set("DSN", randomTestString)

		// in-memory value not changed.
		cfg := NewConfig("testdata/default.yaml")
		require.Equal(t, cfgInMemory.DB.Path, cfg.DB.Path)

		cfg = NewConfig("testdata/sql.yaml")
		require.Equal(t, randomTestString, cfg.DB.Path)
	})
}
