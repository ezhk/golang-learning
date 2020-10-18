package core

import (
	"errors"
	"math"

	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
)

var (
	ErrUnexpectedValue = errors.New("unexpected value")
	ErrSmallValue      = errors.New("too small value")
	ErrNaNValue        = errors.New("undefined value")
)

// Score - based on UCB1 algorythm.
func Score(b *structs.BannerScore, summaryShows int64) (float64, error) {
	fClicks := float64(b.Clicks)
	fShows := float64(b.Shows)
	fSummaryShows := float64(summaryShows)

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
