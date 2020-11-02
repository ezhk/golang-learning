// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/exceptions"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/stretchr/testify/suite"
)

type BannerEventTestSuite struct {
	suite.Suite
	db DatabaseInterface

	firstBanner  *structs.Banner
	secondBanner *structs.Banner
	firstSlot    *structs.Slot
	secondSlot   *structs.Slot
	group        *structs.Group
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(BannerEventTestSuite))
}

func (bs *BannerEventTestSuite) SetupTest() {
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

func (bs *BannerEventTestSuite) TearDownTest() {
	scores, _ := bs.db.ReadBannersPlacements(nil)
	for _, score := range scores {
		_ = bs.db.DeleteBannerPlacement(score.ID)
	}

	_ = bs.db.DeleteBanner(bs.firstBanner.ID)
	_ = bs.db.DeleteBanner(bs.secondBanner.ID)
	_ = bs.db.DeleteSlot(bs.firstSlot.ID)
	_ = bs.db.DeleteSlot(bs.secondSlot.ID)
	_ = bs.db.DeleteGroup(bs.group.ID)
}

func (bs *BannerEventTestSuite) TestBannerEvents() {
	firstScore, err := bs.db.CreateBannerPlacement(bs.firstBanner.ID, bs.firstSlot.ID, bs.group.ID)
	bs.NoError(err)

	secondScore, err := bs.db.CreateBannerPlacement(bs.secondBanner.ID, bs.firstSlot.ID, bs.group.ID)
	bs.NoError(err)

	thirdScore, err := bs.db.CreateBannerPlacement(bs.secondBanner.ID, bs.secondSlot.ID, bs.group.ID)
	bs.NoError(err)

	// First events: exceptions because exists seconds banner with zero shows
	err = bs.db.ProcessBannerEvent(firstScore.ID, "show")
	bs.Error(exceptions.ErrNoChanges, err)
	err = bs.db.ProcessBannerEvent(firstScore.ID, "click")
	bs.Error(exceptions.ErrNoChanges, err)
	err = bs.db.ProcessBannerEvent(firstScore.ID, "show")
	bs.Error(exceptions.ErrNoChanges, err)

	// Second events.
	err = bs.db.ProcessBannerEvent(secondScore.ID, "show")
	bs.NoError(err)
	err = bs.db.ProcessBannerEvent(secondScore.ID, "show")
	bs.NoError(err)
	err = bs.db.ProcessBannerEvent(secondScore.ID, "show")
	bs.NoError(err)

	// Third events.
	err = bs.db.ProcessBannerEvent(thirdScore.ID, "show")
	bs.NoError(err)
	err = bs.db.ProcessBannerEvent(thirdScore.ID, "click")
	bs.NoError(err)
	err = bs.db.ProcessBannerEvent(thirdScore.ID, "show")
	bs.NoError(err)
	err = bs.db.ProcessBannerEvent(thirdScore.ID, "click")
	bs.NoError(err)

	// Update banners scores.
	err = bs.db.RecalculateBannersScore(nil)
	bs.NoError(err)

	// Check values.
	topBanner, err := bs.db.ReadBannerHighestScore(nil)
	bs.NoError(err)
	bs.Equal(thirdScore.ID, topBanner.ID)
	bs.Equal(1.8325546111576978, topBanner.Score)
	bs.Equal(int64(2), topBanner.Shows)
	bs.Equal(int64(2), topBanner.Clicks)

	// Synthetic change leaders.
	bs.db.ProcessBannerEvent(firstScore.ID, "show")
	bs.db.ProcessBannerEvent(firstScore.ID, "click")
	bs.db.ProcessBannerEvent(firstScore.ID, "click")
	bs.db.RecalculateBannersScore(nil)

	topBanner, err = bs.db.ReadBannerHighestScore(nil)
	bs.NoError(err)
	bs.Equal(firstScore.ID, topBanner.ID)
	bs.Equal(2.0929347248663586, topBanner.Score)
	bs.Equal(int64(3), topBanner.Shows)
	bs.Equal(int64(3), topBanner.Clicks)
}
