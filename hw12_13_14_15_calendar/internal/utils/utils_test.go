package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDate(t *testing.T) {
	t.Run("day range", func(t *testing.T) {
		date := time.Date(2020, 9, 17, 12, 30, 14, 0, time.UTC)
		require.Equal(t, time.Date(2020, 9, 17, 0, 0, 0, 0, time.UTC), StartDay(date))
		require.Equal(t, time.Date(2020, 9, 18, 0, 0, 0, 0, time.UTC), EndDay(date))

		date = time.Date(2019, 12, 31, 18, 15, 37, 0, time.UTC)
		require.Equal(t, time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC), StartDay(date))
		require.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), EndDay(date))
	})

	t.Run("week range", func(t *testing.T) {
		date := time.Date(2020, 9, 17, 12, 30, 14, 0, time.UTC)
		require.Equal(t, time.Date(2020, 9, 14, 0, 0, 0, 0, time.UTC), StartWeek(date))
		require.Equal(t, time.Date(2020, 9, 21, 0, 0, 0, 0, time.UTC), EndWeek(date))

		date = time.Date(2020, 1, 3, 18, 15, 37, 0, time.UTC)
		require.Equal(t, time.Date(2019, 12, 30, 0, 0, 0, 0, time.UTC), StartWeek(date))
		require.Equal(t, time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC), EndWeek(date))
	})

	t.Run("month range", func(t *testing.T) {
		date := time.Date(2020, 9, 17, 12, 30, 14, 0, time.UTC)
		require.Equal(t, time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC), StartMonth(date))
		require.Equal(t, time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC), EndMonth(date))

		// last year's month check
		date = time.Date(2019, 12, 28, 11, 12, 13, 0, time.UTC)
		require.Equal(t, time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC), StartMonth(date))
		require.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), EndMonth(date))

		// first year's month check
		date = time.Date(2020, 1, 3, 18, 15, 37, 0, time.UTC)
		require.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), StartMonth(date))
		require.Equal(t, time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), EndMonth(date))
	})
}
