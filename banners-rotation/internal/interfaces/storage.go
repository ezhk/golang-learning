package interfaces

import "github.com/ezhk/golang-learning/banners-rotation/internal/structs"

type Storage interface {
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
