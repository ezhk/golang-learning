package logger

import (
	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"go.uber.org/zap"
)

type Logger struct {
	zap.Logger
}

func NewLogger(config *config.Configuration) *Logger {
	zapConfig := config.ZapConfigBuilder()
	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{*logger}
}

func (l *Logger) Close() error {
	return l.Sync()
}
