package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Check not exists file", func(t *testing.T) {
		err := Copy("not-exist-file", "", 0, 0)
		require.Equal(t, err, ErrFileNotExist)
	})

	t.Run("Too long offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "", int64(math.Exp2(20)), 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})
}
