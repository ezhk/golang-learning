// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/stretchr/testify/suite"
)

type BannerPlacementTestSuite struct {
	suite.Suite
	db *Storage

	firstBanner  *structs.Banner
	secondBanner *structs.Banner
	firstSlot    *structs.Slot
	secondSlot   *structs.Slot
	group        *structs.Group
}

func TestBannerPlacementSuite(t *testing.T) {
	suite.Run(t, new(BannerPlacementTestSuite))
}

func (bs *BannerPlacementTestSuite) SetupTest() {
	cfg := config.NewConfig("testdata/config.yaml")

	db, err := NewStorage(cfg)
	bs.NoError(err)

	// Define Storage.
	bs.db = db

	// Clean exists Groups.
	BannersPlacements, err := bs.db.ReadBannersPlacements(nil)
	bs.NoError(err)
	for _, BannerPlacement := range BannersPlacements {
		err = bs.db.DeleteGroup(BannerPlacement.ID)
		bs.NoError(err)
	}

	// Create test banner.
	firstBanner, err := bs.db.CreateBanner("test first score banner", "first banner description")
	bs.NoError(err)
	bs.firstBanner = &firstBanner

	// Create test banner.
	secondBanner, err := bs.db.CreateBanner("test second score banner", "second banner description")
	bs.NoError(err)
	bs.secondBanner = &secondBanner

	// Create first slot.
	firstSlot, err := bs.db.CreateSlot("test first score slot", "first slot description")
	bs.NoError(err)
	bs.firstSlot = &firstSlot

	// Create second slot.
	secondSlot, err := bs.db.CreateSlot("test second score slot", "second slot description")
	bs.NoError(err)
	bs.secondSlot = &secondSlot

	// Create test banner.
	group, err := bs.db.CreateGroup("test score group", "group description")
	bs.NoError(err)
	bs.group = &group
}

func (bs *BannerPlacementTestSuite) TearDownTest() {
	_ = bs.db.DeleteBanner(bs.firstBanner.ID)
	_ = bs.db.DeleteBanner(bs.secondBanner.ID)
	_ = bs.db.DeleteSlot(bs.firstSlot.ID)
	_ = bs.db.DeleteSlot(bs.secondSlot.ID)
	_ = bs.db.DeleteGroup(bs.group.ID)
}

func (bs *BannerPlacementTestSuite) TestBannerPlacementOperations() {
	// Create new object.
	bScore, err := bs.db.CreateBannerPlacement(bs.firstBanner.ID, bs.firstSlot.ID, bs.group.ID)
	bs.NoError(err)
	defer bs.db.DeleteBannerPlacement(bScore.ID)

	bs.NotNil(bScore)
	bs.Equal(bs.firstBanner.ID, bScore.BannerID)

	// Read and check joined objects.
	scores, err := bs.db.ReadBannersPlacements(nil)
	bs.NoError(err)
	bs.Greater(len(scores), 0)

	bs.Equal(bs.firstBanner.Name, scores[0].Banner.Name)
	bs.Equal(bs.firstSlot.Name, scores[0].Slot.Name)
	bs.Equal(bs.group.Name, scores[0].Group.Name)

	// Update with non exist ID: cause error.
	scores[0].BannerID = scores[0].BannerID + 10
	bScore, err = bs.db.UpdateBannerPlacement(*scores[0])
	bs.Error(err)

	// Valid ID.
	scores[0].BannerID = bs.secondBanner.ID
	bScore, err = bs.db.UpdateBannerPlacement(*scores[0])
	bs.NoError(err)
	bs.Equal(bs.secondBanner.ID, bScore.BannerID)
}

func (bs *BannerPlacementTestSuite) TestReadBannersShows() {
	firstPlacement, err := bs.db.CreateBannerPlacement(bs.firstBanner.ID, bs.firstSlot.ID, bs.group.ID)
	bs.NoError(err)
	defer bs.db.DeleteBannerPlacement(firstPlacement.ID)

	secondPlacement, err := bs.db.CreateBannerPlacement(bs.secondBanner.ID, bs.firstSlot.ID, bs.group.ID)
	bs.NoError(err)
	defer bs.db.DeleteBannerPlacement(secondPlacement.ID)

	scores, err := bs.db.ReadBannersPlacements(nil)
	bs.NoError(err)
	for _, score := range scores {
		score.Shows = 64
		_, err = bs.db.UpdateBannerPlacement(*score)
		bs.NoError(err)
	}

	bScores, err := bs.db.ReadBannersShows(nil)
	bs.NoError(err)
	bs.Equal(1, len(bScores))
	bs.Equal(int64(64*2), bScores[0].SummaryShows)

	_, err = bs.db.CreateBannerPlacement(bs.secondBanner.ID, bs.secondSlot.ID, bs.group.ID)
	bs.NoError(err)

	bScores, err = bs.db.ReadBannersShows(nil)
	bs.NoError(err)
	bs.Equal(2, len(bScores))
	bs.Equal(int64(64*2), bScores[0].SummaryShows+bScores[1].SummaryShows)

	// Test filter operations: don't forget about "GROUP BY"
	bScores, err = bs.db.ReadBannersShows(structs.BannerFilter{"banner_id": bs.firstBanner.ID})
	bs.NoError(err)
	bs.Equal(1, len(bScores))

	bScores, err = bs.db.ReadBannersShows(structs.BannerFilter{"group_id": bs.group.ID})
	bs.NoError(err)
	bs.Equal(2, len(bScores)) // because 2 groups: group_id:FIRST_slot_id, group_id:SECOND_slot_id

	bScores, err = bs.db.ReadBannersShows(structs.BannerFilter{"group_id": bs.group.ID, "slot_id": bs.firstSlot.ID})
	bs.NoError(err)
	bs.Equal(1, len(bScores)) // 1 group: group_id:FIRST_slot_id

	bScores, err = bs.db.ReadBannersShows(structs.BannerFilter{"group_id": bs.group.ID + 1})
	bs.NoError(err)
	bs.Equal(0, len(bScores))
}

func (bs *BannerPlacementTestSuite) TestReadBannerWithHighestScore() {
	// Create banner score and update value.
	firstScore, err := bs.db.CreateBannerPlacement(bs.firstBanner.ID, bs.firstSlot.ID, bs.group.ID)
	bs.NoError(err)
	defer bs.db.DeleteBannerPlacement(firstScore.ID)

	firstScore.Score = 97
	_, err = bs.db.UpdateBannerPlacement(firstScore)
	bs.NoError(err)

	// Create second banner with lowest score value.
	score, err := bs.db.CreateBannerPlacement(bs.secondBanner.ID, bs.firstSlot.ID, bs.group.ID)
	bs.NoError(err)
	defer bs.db.DeleteBannerPlacement(score.ID)

	score.Score = 13
	_, err = bs.db.UpdateBannerPlacement(score)
	bs.NoError(err)

	banner, err := bs.db.ReadBannerHighestScore(nil)
	bs.NoError(err)
	bs.Equal(bs.firstBanner.ID, banner.BannerID)
}
