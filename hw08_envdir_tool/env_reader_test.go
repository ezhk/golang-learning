package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadFileValue(t *testing.T) {
	t.Run("not exist file", func(t *testing.T) {
		_, err := readFileValue("testdata/env/not-exist")
		require.True(t, os.IsNotExist(err))
	})

	t.Run("empty file", func(t *testing.T) {
		val, err := readFileValue("testdata/env/UNSET")
		require.Nil(t, err)
		require.Equal(t, "", val)
	})

	t.Run("heavy variable", func(t *testing.T) {
		sample := fmt.Sprintf("   foo\nwith new line")

		val, err := readFileValue("testdata/env/FOO")
		require.Nil(t, err)
		require.Equal(t, sample, val)
	})
}

func TestReadDir(t *testing.T) {
	t.Run("not exist dir", func(t *testing.T) {
		_, err := ReadDir("testdata/env-not-exist")
		err, ok := err.(*os.PathError)

		require.True(t, ok)
		require.NotNil(t, err)
	})

	t.Run("complex parse check", func(t *testing.T) {
		envMap, err := ReadDir("testdata/env")
		require.Nil(t, err)

		require.Equal(t, envMap, Environment{"UNSET": "", "BAR": "bar", "HELLO": "\"hello\"", "FOO": "   foo\nwith new line"})
	})
}
