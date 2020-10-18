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
