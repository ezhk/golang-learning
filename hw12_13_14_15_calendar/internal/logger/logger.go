package logger

import (
	"log"

	config "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
)

type Logger struct {
	zap.SugaredLogger
	Sync func() error
}

func NewLogger(cfg *config.Configuration) *Logger {
	zapConfig := cfg.ZapConfigBuilder()
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	sugar := logger.Sugar()

	return &Logger{
		*sugar,
		logger.Sync,
	}
}

func (l *Logger) Close() error {
	return l.Sync()
}
