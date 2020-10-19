// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/stretchr/testify/suite"
)

type BannerScoreTestSuite struct {
	suite.Suite
	db *Storage

	firstBanner  *structs.Banner
	secondBanner *structs.Banner
	firstSlot    *structs.Slot
	secondSlot   *structs.Slot
	group        *structs.Group
}

func TestBannerScoreSuite(t *testing.T) {
	suite.Run(t, new(BannerScoreTestSuite))
}

func (bs *BannerScoreTestSuite) SetupTest() {
	cfg := config.NewConfig("testdata/config.yaml")

	db, err := NewStorage(cfg)
	bs.NoError(err)

	// Define Storage.
	bs.db = db

	// Clean exists Groups.
	bannersScores, err := bs.db.ReadBannersScores()
	bs.NoError(err)
	for _, bannerScore := range bannersScores {
		err = bs.db.DeleteGroup(bannerScore.ID)
		bs.NoError(err)
	}

	// Create test banner.
	firstBanner, err := bs.db.CreateBanner("test first score banner")
	bs.NoError(err)
	bs.firstBanner = &firstBanner

	// Create test banner.
	secondBanner, err := bs.db.CreateBanner("test second score banner")
	bs.NoError(err)
	bs.secondBanner = &secondBanner

	// Create first slot.
	firstSlot, err := bs.db.CreateSlot("test first score slot")
	bs.NoError(err)
	bs.firstSlot = &firstSlot

	// Create second slot.
	secondSlot, err := bs.db.CreateSlot("test second score slot")
	bs.NoError(err)
	bs.secondSlot = &secondSlot

	// Create test banner.
	group, err := bs.db.CreateGroup("test score group")
	bs.NoError(err)
	bs.group = &group
}

func (bs *BannerScoreTestSuite) TearDownTest() {
	scores, _ := bs.db.ReadBannersScores()
	for _, score := range scores {
		_ = bs.db.DeleteBannerScore(score.ID)
	}

	_ = bs.db.DeleteBanner(bs.firstBanner.ID)
	_ = bs.db.DeleteBanner(bs.secondBanner.ID)
	_ = bs.db.DeleteSlot(bs.firstSlot.ID)
	_ = bs.db.DeleteSlot(bs.secondSlot.ID)
	_ = bs.db.DeleteGroup(bs.group.ID)
}

func (bs *BannerScoreTestSuite) TestBannerScoreOperations() {
	// // Delete operations.
	// scores, err := bs.db.ReadBannersScores()
	// bs.NoError(err)
	// for _, score := range scores {
	// 	err = bs.db.DeleteBannerScore(score.ID)
	// 	bs.NoError(err)
	// }

	// Create new object.
	bScore, err := bs.db.CreateBannerScore(bs.firstBanner.ID, bs.firstSlot.ID, bs.group.ID)
	bs.NoError(err)
	bs.NotNil(bScore)
	bs.Equal(bs.firstBanner.ID, bScore.BannerID)

	// Read and check joined objects.
	scores, err := bs.db.ReadBannersScores()
	bs.NoError(err)
	bs.Equal(1, len(scores))
	bs.Equal(bs.firstBanner.Name, scores[0].Banner.Name)
	bs.Equal(bs.firstSlot.Name, scores[0].Slot.Name)
	bs.Equal(bs.group.Name, scores[0].Group.Name)

	// Update with non exist ID: cause error.
	scores[0].BannerID = scores[0].BannerID + 10
	bScore, err = bs.db.UpdateBannerScore(*scores[0])
	bs.Error(err)

	// Valid ID.
	scores[0].BannerID = bs.secondBanner.ID
	bScore, err = bs.db.UpdateBannerScore(*scores[0])
	bs.NoError(err)
	bs.Equal(bs.secondBanner.ID, bScore.BannerID)

	// Delete operations.
	scores, err = bs.db.ReadBannersScores()
	bs.NoError(err)
	bs.Equal(1, len(scores))
	for _, score := range scores {
		err = bs.db.DeleteBannerScore(score.ID)
		bs.NoError(err)
	}
}

func (bs *BannerScoreTestSuite) TestReadBannersShows() {
	// // Delete operations.
	// scores, err := bs.db.ReadBannersScores()
	// bs.NoError(err)
	// for _, score := range scores {
	// 	err = bs.db.DeleteBannerScore(score.ID)
	// 	bs.NoError(err)
	// }

	_, err := bs.db.CreateBannerScore(bs.firstBanner.ID, bs.firstSlot.ID, bs.group.ID)
	bs.NoError(err)

	_, err = bs.db.CreateBannerScore(bs.secondBanner.ID, bs.firstSlot.ID, bs.group.ID)
	bs.NoError(err)

	scores, err := bs.db.ReadBannersScores()
	bs.NoError(err)
	for _, score := range scores {
		score.Shows = 64
		_, err = bs.db.UpdateBannerScore(*score)
		bs.NoError(err)
	}

	bScores, err := bs.db.ReadBannersShows()
	bs.NoError(err)
	bs.Equal(1, len(bScores))
	bs.Equal(int64(64*2), bScores[0].SummaryShows)

	_, err = bs.db.CreateBannerScore(bs.secondBanner.ID, bs.secondSlot.ID, bs.group.ID)
	bs.NoError(err)

	bScores, err = bs.db.ReadBannersShows()
	bs.NoError(err)
	bs.Equal(2, len(bScores))
	bs.Equal(int64(64*2), bScores[0].SummaryShows+bScores[1].SummaryShows)
}
