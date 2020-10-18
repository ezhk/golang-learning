package logger

import (
	"log"

	config "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
)

type Logger struct {
	zap.Logger
}

func NewLogger(cfg *config.Configuration) *Logger {
	zapConfig := cfg.ZapConfigBuilder()
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return &Logger{*logger}
}

func (l *Logger) Close() error {
	return l.Sync()
}
