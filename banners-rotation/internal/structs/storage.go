package structs

import (
	"time"
)

type Slot struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string `gorm:"unique"`
	Description string
}

type Group struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string `gorm:"unique"`
	Description string
}

type Banner struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string `gorm:"unique"`
	Description string
}

type BannerScore struct {
	// gorm.Model
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time

	BannerID uint64
	Banner   Banner `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	SlotID uint64
	Slot   Slot `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	GroupID uint64
	Group   Group `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Shows  int64   `gorm:"default:0"`
	Clicks int64   `gorm:"default:0"`
	Score  float64 `gorm:"default:0"` // will be computed in background task
}

// SummaryBannersShows using for UCB1 summary show value.
type SummaryBannersShows struct {
	SlotID       uint64
	GroupID      uint64
	SummaryShows int64
}
