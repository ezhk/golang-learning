// +build !race

package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExtendedTelnetClient(t *testing.T) {
	t.Run("extended", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer l.Close()

		var wg sync.WaitGroup
		wg.Add(2)

		// client part
		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer client.Close()

			in.WriteString("HELLO\n")
			go client.Sender()
			go client.Receiver()

			// wait 1 sec for receive message from server
			time.Sleep(1 * time.Second)
			require.Equal(t, "WORLD\n", out.String())
		}()

		// server part
		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer conn.Close()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)

			// client received empty message if input message is incorrect
			require.Equal(t, "HELLO\n", string(request)[:n])

			n, err = conn.Write([]byte("WORLD\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}
