package logger

import (
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Logger struct {
	zap.SugaredLogger
	Sync func() error
}

func NewLogger() *Logger {
	config := prepareConfig()

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	sugar := logger.Sugar()

	return &Logger{
		*sugar,
		logger.Sync,
	}
}

func prepareConfig() zap.Config {
	config := zap.NewProductionConfig()
	if v := viper.GetString("logger.encoding"); len(v) > 0 {
		config.Encoding = v
	}

	if v := viper.GetString("logger.level"); len(v) > 0 {
		switch v {
		case "debug":
			config.Level.SetLevel(zap.DebugLevel)
		case "info":
			config.Level.SetLevel(zap.InfoLevel)
		case "warning":
			config.Level.SetLevel(zap.WarnLevel)
		case "error":
			config.Level.SetLevel(zap.ErrorLevel)
		}
	}

	if v := viper.GetStringSlice("logger.output"); len(v) > 0 {
		config.OutputPaths = v
	}

	if v := viper.GetStringSlice("logger.error"); len(v) > 0 {
		config.ErrorOutputPaths = v
	}

	return config
}

func (l *Logger) Close() error {
	return l.Sync()
}
