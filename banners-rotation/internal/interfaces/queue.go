package interfaces

import (
	"github.com/ezhk/golang-learning/banners-rotation/internal/logger"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
)

type Queue interface {
	ProduceEvent(structs.QueueEvent) error

	RunConsumer(Storage, *logger.Logger) error
	StopConsumer() error
}
