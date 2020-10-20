package storage

import (
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"gorm.io/gorm/clause"
)

func (s *Storage) CreateBannerScore(bannerID uint64, slotID uint64, groupID uint64) (structs.BannerScore, error) {
	bannerScore := structs.BannerScore{
		BannerID: bannerID,
		SlotID:   slotID,
		GroupID:  groupID,
	}
	result := s.db.Create(&bannerScore)

	return bannerScore, result.Error
}

func (s *Storage) ReadBannersScores() ([]*structs.BannerScore, error) {
	bannersScores := []*structs.BannerScore{}
	result := s.db.Preload(clause.Associations).Find(&bannersScores)

	return bannersScores, result.Error
}

func (s *Storage) ReadBannersShows() ([]*structs.SummaryBannersShows, error) {
	bannersShows := []*structs.SummaryBannersShows{}
	result := s.db.Preload(clause.Associations).Table("banner_scores").Select("slot_id, group_id, sum(shows) as summary_shows").Group("slot_id").Group("group_id").Find(&bannersShows)

	return bannersShows, result.Error
}

func (s *Storage) UpdateBannerScore(b structs.BannerScore) (structs.BannerScore, error) {
	bannerScore := structs.BannerScore{ID: b.ID}
	result := s.db.Model(&bannerScore).Updates(b)

	return bannerScore, result.Error
}

func (s *Storage) DeleteBannerScore(id uint64) error {
	result := s.db.Delete(&structs.BannerScore{}, id)

	return result.Error
}
