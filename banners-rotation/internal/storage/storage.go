package storage

import (
	"fmt"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseInterface interface {
	CreateBanner(string, string) (structs.Banner, error)
	ReadBanners() ([]*structs.Banner, error)
	UpdateBanner(structs.Banner) (structs.Banner, error)
	DeleteBanner(uint64) error

	CreateSlot(string, string) (structs.Slot, error)
	ReadSlots() ([]*structs.Slot, error)
	UpdateSlot(structs.Slot) (structs.Slot, error)
	DeleteSlot(uint64) error

	CreateGroup(string, string) (structs.Group, error)
	ReadGroups() ([]*structs.Group, error)
	UpdateGroup(structs.Group) (structs.Group, error)
	DeleteGroup(uint64) error

	CreateBannerPlacement(uint64, uint64, uint64) (structs.BannerPlacement, error)
	ReadBannerHighestScore(structs.BannerFilter) (structs.BannerPlacement, error)
	ReadBannersPlacements(structs.BannerFilter) ([]*structs.BannerPlacement, error)
	ReadBannersShows(structs.BannerFilter) ([]*structs.SummaryBannersShows, error)
	UpdateBannerPlacement(structs.BannerPlacement) (structs.BannerPlacement, error)
	DeleteBannerPlacement(uint64) error

	ProcessBannerEvent(uint64, string) error
	RecalculateBannersScore(structs.BannerFilter) error
}

type Storage struct {
	db *gorm.DB
}

func NewStorage(cfg *config.Configuration) (DatabaseInterface, error) {
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
