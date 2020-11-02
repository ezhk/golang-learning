// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDatabase(t *testing.T) {
	t.Run("new database call: check migrations apply", func(t *testing.T) {
		cfg := config.NewConfig("testdata/config.yaml")
		db, err := NewStorage(cfg)

		require.NoError(t, err)
		require.NotNil(t, db)
	})
}

func TestFilter(t *testing.T) {
	t.Run("filter test", func(t *testing.T) {
		cfg := config.NewConfig("testdata/config.yaml")
		db, err := gorm.Open(postgres.Open(cfg.GetDatabasePath()), &gorm.Config{})
		require.NoError(t, err)

		err = FilterMap(db, structs.BannerFilter{"group_id": 123})
		require.NoError(t, err)
	})
}
