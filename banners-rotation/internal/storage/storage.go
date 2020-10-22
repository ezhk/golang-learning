package storage

import (
	"fmt"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

// Generate base methods for tables.
//go:generate go run ./generate-table-methods/... -table banner -file generated_banner.go
//go:generate go run ./generate-table-methods/... -table slot -file generated_slot.go
//go:generate go run ./generate-table-methods/... -table group -file generated_group.go

// Generate base tests for table methods.
//go:generate go run ./generate-table-tests/... -table banner -file generated_banner_test.go
//go:generate go run ./generate-table-tests/... -table slot -file generated_slot_test.go
//go:generate go run ./generate-table-tests/... -table group -file generated_group_test.go

func NewStorage(cfg *config.Configuration) (*Storage, error) {
	// Create psql connection.
	db, err := gorm.Open(postgres.Open(cfg.GetDatabasePath()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Apply migration on create storage object.
	err = db.AutoMigrate(&structs.Slot{}, &structs.Group{}, &structs.Banner{}, &structs.BannerPlacement{})
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func FilterMap(tx *gorm.DB, filter structs.BannerFilter) error {
	for k, v := range filter {
		tx = tx.Where(fmt.Sprintf("%s = ?", k), v)
	}

	return nil
}
