package storage

import (
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"gorm.io/gorm/clause"
)

func (s *Storage) CreateBannerPlacement(bannerID uint64, slotID uint64, groupID uint64) (structs.BannerPlacement, error) {
	bannerPlacement := structs.BannerPlacement{
		BannerID: bannerID,
		SlotID:   slotID,
		GroupID:  groupID,
	}
	result := s.db.Create(&bannerPlacement)

	return bannerPlacement, result.Error
}

func (s *Storage) ReadBannerHighestScore(filter structs.BannerFilter) (structs.BannerPlacement, error) {
	bannerPlacement := structs.BannerPlacement{}
	query := s.db.Preload(clause.Associations).Order("score desc")
	if err := FilterMap(query, filter); err != nil {
		return structs.BannerPlacement{}, err
	}
	result := query.First(&bannerPlacement)

	return bannerPlacement, result.Error
}

func (s *Storage) ReadBannersPlacements(filter structs.BannerFilter) ([]*structs.BannerPlacement, error) {
	manyBannersPlacements := []*structs.BannerPlacement{}
	query := s.db.Preload(clause.Associations)
	if err := FilterMap(query, filter); err != nil {
		return manyBannersPlacements, err
	}
	result := query.Find(&manyBannersPlacements)

	return manyBannersPlacements, result.Error
}

func (s *Storage) ReadBannersShows(filter structs.BannerFilter) ([]*structs.SummaryBannersShows, error) {
	manyBannersShows := []*structs.SummaryBannersShows{}

	query := s.db.Preload(clause.Associations).Table("banner_placements")
	query = query.Select("slot_id, group_id, sum(shows) as summary_shows").Group("slot_id").Group("group_id")

	if err := FilterMap(query, filter); err != nil {
		return manyBannersShows, err
	}

	result := query.Find(&manyBannersShows)

	return manyBannersShows, result.Error
}

func (s *Storage) UpdateBannerPlacement(b structs.BannerPlacement) (structs.BannerPlacement, error) {
	bannerPlacement := structs.BannerPlacement{ID: b.ID}
	result := s.db.Model(&bannerPlacement).Updates(b)

	return bannerPlacement, result.Error
}

func (s *Storage) DeleteBannerPlacement(id uint64) error {
	result := s.db.Delete(&structs.BannerPlacement{}, id)

	return result.Error
}
