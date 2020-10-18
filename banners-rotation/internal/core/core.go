package core

import (
	"errors"
	"math"
)

var (
	ErrUnexpectedValue = errors.New("unexpected value")
	ErrSmallValue      = errors.New("too small value")
	ErrNaNValue        = errors.New("undefined value")
)

// Score - based on UCB1 algorythm.
func Score(clicks int64, bannerShows int64, summaryBannersShows int64) (float64, error) {
	fClicks := float64(clicks)
	fShows := float64(bannerShows)
	fSummaryShows := float64(summaryBannersShows)

	if fShows < 1 {
		return 0, ErrSmallValue
	}

	if fClicks > fShows || fShows > fSummaryShows {
		return 0, ErrUnexpectedValue
	}

	averageReward := fClicks / fShows
	val := averageReward + math.Sqrt(2*math.Log(fSummaryShows)/fShows)
	if math.IsNaN(val) {
		return 0, ErrNaNValue
	}

	return val, nil
}
