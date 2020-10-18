// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/stretchr/testify/suite"
)

type BannersTestSuite struct {
	suite.Suite
	db *Storage
}

func TestBannersSuite(t *testing.T) {
	suite.Run(t, new(BannersTestSuite))
}

func (s *BannersTestSuite) SetupTest() {
	cfg := config.NewConfig("testdata/config.yaml")

	db, err := NewStorage(cfg)
	s.NoError(err)

	// Define storage.
	s.db = db

	// Clean exists Banners.
	banners, err := s.db.ReadBanners()
	s.NoError(err)
	for _, Banner := range banners {
		err = s.db.DeleteBanner(Banner.ID)
		s.NoError(err)
	}
}

func (s *BannersTestSuite) TestBannerOperations() {
	// Create new Banner.
	banner, err := s.db.CreateBanner("test banner")
	s.NoError(err)
	s.Equal("test banner", banner.Name)

	banner.Name = "updated test banner"
	updatedBanner, err := s.db.UpdateBanner(banner)
	s.NoError(err)
	s.Equal("updated test banner", updatedBanner.Name)

	banners, err := s.db.ReadBanners()
	s.NoError(err)
	s.Greater(len(banners), 0)
	s.Equal("updated test banner", banners[0].Name)

	// Call "duplicate key value violates unique constraint".
	_, err = s.db.CreateBanner("updated test banner")
	s.Error(err)

	err = s.db.DeleteBanner(banners[0].ID)
	s.NoError(err)
}