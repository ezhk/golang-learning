// +build integration

package broker

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/structs"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	t.Run("message event", func(t *testing.T) {
		event := structs.Event{
			UserID:   1,
			Title:    "Base titles",
			Content:  "Random content",
			DateFrom: time.Date(2020, 1, 2, 12, 4, 37, 0, time.UTC),
			DateTo:   time.Date(2020, 1, 3, 9, 15, 0, 0, time.UTC),
			Notified: false,
		}

		cfg := config.NewConfig("testdata/config.yaml")
		producer := NewProducer(cfg)
		consumer := NewConsumer(cfg)

		err := producer.Init()
		require.NoError(t, err)
		defer producer.Close()

		err = consumer.Init()
		require.NoError(t, err)
		defer consumer.Close()

		// Must runned before producer.
		msgCh, err := consumer.Consume()
		require.NoError(t, err)

		err = producer.Publish(event)
		require.NoError(t, err)

		// Read only one message
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		var Message amqp.Delivery
		select {
		case Message = <-msgCh:
		case <-ctx.Done():
			require.Fail(t, "no read messages for 3 seconds")
		}
		err = Message.Ack(false)
		require.NoError(t, err)

		// Conert message body into base struct and validate.
		var recvEvent structs.Event
		err = json.Unmarshal(Message.Body, &recvEvent)
		require.NoError(t, err)
		require.Equal(t, event, recvEvent)
	})
}
