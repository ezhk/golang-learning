package structs

import "gorm.io/gorm"

type Slot struct {
	gorm.Model
	Name string
}

type Group struct {
	gorm.Model
	Name string
}

type Banner struct {
	gorm.Model
	Name string
}

type BannerScore struct {
	gorm.Model

	Banner Banner
	Slot   Slot
	Group  Group

	Shows  int64
	Clicks int64
	Score  float64 // will be computed in background task
}
