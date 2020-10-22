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
		cfg := config.NewConfig("testdata/config.yaml")
		q, err := NewQueue(cfg)
		require.NoError(t, err)

		err = q.ProduceEvent(structs.QueueEvent{PlacementID: 1, EventType: "click"})
		require.NoError(t, err)

		storage, err := storage.NewStorage(cfg)
		require.NoError(t, err)

		err = q.RunConsumer(storage, nil)
		require.NoError(t, err)

		<-q.Conn.StopAllConsuming()
	})
}
