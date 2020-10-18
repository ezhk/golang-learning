package storage

import "github.com/ezhk/golang-learning/banners-rotation/internal/structs"

func (s *Storage) CreateBanner(bannerName string) (structs.Banner, error) {
	banner := structs.Banner{Name: bannerName}
	result := s.db.Create(&banner)

	return banner, result.Error
}

func (s *Storage) ReadBanners() ([]*structs.Banner, error) {
	banners := []*structs.Banner{}
	result := s.db.Find(&banners)

	return banners, result.Error
}

func (s *Storage) UpdateBanner(b structs.Banner) (structs.Banner, error) {
	banner := structs.Banner{ID: b.ID}
	result := s.db.Model(&banner).Updates(b)

	return banner, result.Error
}

func (s *Storage) DeleteBanner(id uint) error {
	result := s.db.Delete(&structs.Banner{}, id)

	return result.Error
}