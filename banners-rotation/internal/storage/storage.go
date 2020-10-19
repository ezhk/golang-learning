package storage

import (
	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

// Generate base methods for tables.
//go:generate go run ./generate-table-methods/... -table banner -file banner.go
//go:generate go run ./generate-table-methods/... -table slot -file slot.go
//go:generate go run ./generate-table-methods/... -table group -file group.go

// Generate base tests for table methods.
//go:generate go run ./generate-table-tests/... -table banner -file banner_test.go
//go:generate go run ./generate-table-tests/... -table slot -file slot_test.go
//go:generate go run ./generate-table-tests/... -table group -file group_test.go

func NewStorage(cfg *config.Configuration) (*Storage, error) {
	// Create psql connection.
	db, err := gorm.Open(postgres.Open(cfg.GetDatabasePath()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Apply migration on create storage object.
	err = db.AutoMigrate(&structs.Slot{}, &structs.Group{}, &structs.Banner{}, &structs.BannerScore{})
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}
