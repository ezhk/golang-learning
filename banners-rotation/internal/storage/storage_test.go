// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	t.Run("new database call: check migrations apply", func(t *testing.T) {
		cfg := config.NewConfig("testdata/config.yaml")
		db, err := NewStorage(cfg)

		require.NoError(t, err)
		require.NotNil(t, db)
	})
}
