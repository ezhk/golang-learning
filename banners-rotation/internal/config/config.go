package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Configuration struct {
	Logger   LoggerParams   `yaml:"logger"`
	Database DatabaseParams `yaml:"database"`
	Server   ServerParams   `yaml:"server"`
}

type LoggerParams struct {
	Output   []string `yaml:"output"`
	Error    []string `yaml:"error"`
	Level    string   `yaml:"level"`
	Encoding string   `yaml:"encoding"`
}

type DatabaseParams struct {
	Path string `yaml:"path"`
}

type ServerParams struct {
	HTTP string `yaml:"http"`
	GRPC string `yaml:"grpc"`
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
	// Read environment DSN value.
	if v := viper.Get("DSN"); v != nil {
		cfg.Database.Path = v.(string)
	}
}

func (cfg *Configuration) GetDatabasePath() string {
	return cfg.Database.Path
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
