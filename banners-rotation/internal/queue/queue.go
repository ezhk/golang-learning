package queue

import (
	"encoding/json"
	"time"

	"github.com/adjust/rmq/v3"
	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/logger"
	"github.com/ezhk/golang-learning/banners-rotation/internal/storage"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
)

type Queue struct {
	Conn    rmq.Connection
	Events  rmq.Queue
	Storage storage.DatabaseInterface
	Logger  *logger.Logger
}

func NewQueue(cfg *config.Configuration) (*Queue, error) {
	// Create redis connection.
	conn, err := rmq.OpenConnection("queue-events", "tcp", cfg.Queue.Path, 1, nil)
	if err != nil {
		return nil, err
	}

	eventsQueue, err := conn.OpenQueue("events")
	if err != nil {
		return nil, err
	}

	return &Queue{
		Conn:   conn,
		Events: eventsQueue,
	}, nil
}

func (q *Queue) ProduceEvent(e structs.QueueEvent) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	if err := q.Events.PublishBytes(body); err != nil {
		return err
	}

	return nil
}

func (q *Queue) RunConsumer(databaseI storage.DatabaseInterface, logger *logger.Logger) error {
	q.Storage = databaseI
	q.Logger = logger

	if err := q.Events.StartConsuming(10, 100*time.Millisecond); err != nil {
		return err
	}
	if _, err := q.Events.AddConsumer("consumer", q); err != nil {
		return err
	}

	return nil
}

// Internal consumert method.
func (q *Queue) Consume(delivery rmq.Delivery) {
	payload := delivery.Payload()
	q.Logger.Sugar().Infof("received message: %+v", payload)

	var event structs.QueueEvent
	err := json.Unmarshal([]byte(payload), &event)
	if err != nil {
		return
	}

	err = q.Storage.ProcessBannerEvent(event.PlacementID, event.EventType)
	if err != nil {
		q.Logger.Sugar().Errorf("process message error: %s", err)
		if err := delivery.Reject(); err != nil {
			q.Logger.Sugar().Errorf("reject message error: %s", err)
		}

		return
	}

	if err := delivery.Ack(); err != nil {
		q.Logger.Sugar().Errorf("acknowledge message error: %s", err)
	}
}
