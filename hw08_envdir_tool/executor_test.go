package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("not exist file", func(t *testing.T) {
		exitCode := RunCmd([]string{"not-exist-command"}, nil)
		require.EqualValues(t, 1, exitCode)
	})

	t.Run("check output", func(t *testing.T) {
		r, w, err := os.Pipe()
		require.Nil(t, err)

		os.Stdout = w
		exitCode := RunCmd([]string{"echo", "-n", "test"}, nil)
		w.Close()

		require.EqualValues(t, 0, exitCode)
		out, err := ioutil.ReadAll(r)
		require.Nil(t, err)

		require.Equal(t, "test", string(out))
	})
}
