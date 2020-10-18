package core

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/stretchr/testify/require"
)

type statTest struct {
	shows        int64
	clicks       int64
	summaryShows int64

	expected float64
	err      error
}

func TestCore(t *testing.T) {
	t.Run("UCB1 values table tests", func(t *testing.T) {
		for _, tt := range [...]statTest{
			{
				shows:        0,
				clicks:       0,
				summaryShows: 0,
				expected:     0,
				err:          ErrSmallValue,
			},
			{
				shows:        -10,
				clicks:       2,
				summaryShows: 50,
				expected:     0,
				err:          ErrSmallValue,
			},
			{
				shows:        1,
				clicks:       20,
				summaryShows: 50,
				expected:     0,
				err:          ErrUnexpectedValue,
			},
			{
				shows:        256,
				clicks:       32,
				summaryShows: 128,
				expected:     0,
				err:          ErrUnexpectedValue,
			},
			{
				shows:        10,
				clicks:       2,
				summaryShows: 50,
				expected:     1.0845363763495706,
				err:          nil,
			},
			{
				shows:        100,
				clicks:       23,
				summaryShows: 105,
				expected:     0.535088850997788,
				err:          nil,
			},
			{
				shows:        2,
				clicks:       0,
				summaryShows: 3,
				expected:     1.0481470739682048,
				err:          nil,
			},
			{
				shows:        1,
				clicks:       0,
				summaryShows: 3,
				expected:     1.482303807367511,
				err:          nil,
			},
		} {
			scoreVal, err := Score(&structs.BannerScore{
				Shows:  tt.shows,
				Clicks: tt.clicks,
			}, tt.summaryShows)

			require.Equal(t, tt.expected, scoreVal)
			require.Equal(t, tt.err, err)
		}
	})
}
