// +build integration

package queue

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/storage"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("generate event", func(t *testing.T) {
		cfg := &config.Configuration{
			Queue:    config.QueueParams{Path: "localhost:6379"},
			Database: config.DatabaseParams{Path: "user=postgres password=postgres dbname=postgres sslmode=disable host=localhost port=5432"},
		}

		q, err := NewQueue(cfg)
		require.NoError(t, err)

		err = q.ProduceEvent(structs.QueueEvent{PlacementID: 1, EventType: "click"})
		require.NoError(t, err)

		storage, err := storage.NewStorage(cfg)
		require.NoError(t, err)

		err = q.RunConsumer(storage, nil)
		require.NoError(t, err)

		err = q.StopConsumer()
		require.NoError(t, err)
	})
}
