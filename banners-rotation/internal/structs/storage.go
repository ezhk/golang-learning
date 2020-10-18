package structs

import (
	"time"
)

type Slot struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string `gorm:"unique"`
}

type Group struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string `gorm:"unique"`
}

type Banner struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string `gorm:"unique"`
}

type BannerScore struct {
	// gorm.Model
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	BannerID uint
	Banner   Banner `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	SlotID uint
	Slot   Slot `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	GroupID uint
	Group   Group `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Shows  int64   `gorm:"default:0"`
	Clicks int64   `gorm:"default:0"`
	Score  float64 `gorm:"default:0"` // will be computed in background task
}
