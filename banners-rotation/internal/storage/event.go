package storage

import (
	"fmt"

	"github.com/ezhk/golang-learning/banners-rotation/internal/core"
	"github.com/ezhk/golang-learning/banners-rotation/internal/exceptions"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
)

func (s *Storage) ProcessBannerEvent(id uint64, event string) error {
	manyBannersPlacements, err := s.ReadBannersPlacements(structs.BannerFilter{"id": id})
	if err != nil {
		return err
	}
	if len(manyBannersPlacements) < 1 {
		return exceptions.ErrObjectNotExist
	}

	bannerPlacement := manyBannersPlacements[0]
	switch event {
	case "show":
		bannerPlacement.Shows++
	case "click":
		bannerPlacement.Clicks++
	}

	_, err = s.UpdateBannerPlacement(*bannerPlacement)
	if err != nil {
		return err
	}

	// TODO: too heavy logic, move it to background
	err = s.RecalculateBannersScore(structs.BannerFilter{"slot_id": bannerPlacement.SlotID, "group_id": bannerPlacement.GroupID})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RecalculateBannersScore(filter structs.BannerFilter) error {
	// Prepare summary stats in memory.
	showsInMemory := make(map[string]int64)
	shows, err := s.ReadBannersShows(filter)
	if err != nil {
		return err
	}

	for _, rec := range shows {
		key := fmt.Sprintf("%d-%d", rec.SlotID, rec.GroupID)
		showsInMemory[key] = rec.SummaryShows
	}

	// Update banners stats.
	manyBannersPlacements, err := s.ReadBannersPlacements(filter)
	if err != nil {
		return err
	}

	for _, bannerPlacement := range manyBannersPlacements {
		key := fmt.Sprintf("%d-%d", bannerPlacement.SlotID, bannerPlacement.GroupID)
		sum, ok := showsInMemory[key]
		if !ok {
			return exceptions.ErrKeyNotExist
		}

		val, err := core.Score(bannerPlacement.Clicks, bannerPlacement.Shows, sum)
		if err != nil {
			// Skip current store value, but process another banners.
			continue
		}

		bannerPlacement.Score = val
		_, err = s.UpdateBannerPlacement(*bannerPlacement)
		if err != nil {
			return err
		}
	}

	return nil
}
