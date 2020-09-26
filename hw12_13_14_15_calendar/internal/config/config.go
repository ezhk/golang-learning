package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Configuration struct {
	Logger LoggerParams   `yaml:"logger"`
	DB     DatabaseParams `yaml:"db"`
	Server ServerParams   `yaml:"server"`
	Test   TestParams     `yaml:"test"`
}

type LoggerParams struct {
	Output   []string `yaml:"output"`
	Error    []string `yaml:"error"`
	Level    string   `yaml:"level"`
	Encoding string   `yaml:"encoding"`
}

type DatabaseParams struct {
	Path string `yaml:"path"`
	Type string `yaml:"type"`
}

type ServerParams struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type TestParams struct {
	DoNotRemoveMe int `yaml:"doNotRemoveMe"`
}

func NewConfig(configPath string) *Configuration {
	cfg := Configuration{}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err == nil {
		_ = viper.Unmarshal(&cfg)
	}
	cfg.updateDatabaseConfig()

	return &cfg
}

func (cfg *Configuration) updateDatabaseConfig() {
	switch cfg.DB.Type {
	case "sql":
		// Read environment directory.
		if v := viper.Get("DSN"); v != nil {
			cfg.DB.Path = v.(string)
		}
	case "in-memory":
	}
}

func (cfg *Configuration) ZapConfigBuilder() zap.Config {
	zapConfig := zap.NewProductionConfig()

	if v := cfg.Logger.Encoding; v != "" {
		zapConfig.Encoding = v
	}

	if v := cfg.Logger.Level; v != "" {
		switch v {
		case "error":
			zapConfig.Level.SetLevel(zap.ErrorLevel)
		case "warning":
			zapConfig.Level.SetLevel(zap.WarnLevel)
		case "info":
			zapConfig.Level.SetLevel(zap.InfoLevel)
		default:
			zapConfig.Level.SetLevel(zap.DebugLevel)
		}
	}

	if v := cfg.Logger.Output; len(v) > 0 {
		zapConfig.OutputPaths = v
	}

	if v := cfg.Logger.Error; len(v) > 0 {
		zapConfig.ErrorOutputPaths = v
	}

	return zapConfig
}
