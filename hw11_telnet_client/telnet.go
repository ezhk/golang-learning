package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"time"
)

var (
	ErrEOF           = errors.New("EOF received")
	ErrContextClosed = errors.New("context was closed")
	ErrPeerClosed    = errors.New("connection was closed by peer")
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error

	Sender() error
	Receiver() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &Telnet{Address: address, Timeout: timeout, Reader: in, Writer: out, Ctx: ctx, Cancel: cancel}
}

type Telnet struct {
	Address string
	Timeout time.Duration
	Reader  io.ReadCloser
	Writer  io.Writer

	Ctx    context.Context
	Cancel context.CancelFunc

	Conn net.Conn
}

func (t *Telnet) Connect() error {
	d := net.Dialer{Timeout: t.Timeout}
	conn, err := d.DialContext(t.Ctx, "tcp", t.Address)
	if err != nil {
		return err
	}
	t.Conn = conn

	return nil
}

func (t *Telnet) Close() error {
	defer t.Cancel()

	return t.Conn.Close()
}

func (t *Telnet) Send() error {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(t.Reader)
	if err != nil {
		return err
	}

	_, err = t.Conn.Write(buf.Bytes())
	if err != nil {
		return ErrPeerClosed
	}

	return nil
}

func (t *Telnet) Receive() error {
	data, err := ioutil.ReadAll(t.Conn)
	if err != nil {
		return err
	}

	_, err = t.Writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (t *Telnet) Sender() error {
	scanner := bufio.NewScanner(t.Reader)

	scannerDataCh := make(chan []byte)
	errorCh := make(chan error)

	go func() {
		// scanner.Scan() is lockable, moved to goroutine
		for scanner.Scan() {
			select {
			case <-t.Ctx.Done():
				return
			case scannerDataCh <- scanner.Bytes():
			}
		}

		// close context, when Scan ended
		errorCh <- ErrEOF
	}()

	var payload []byte
	for {
		select {
		case <-t.Ctx.Done():
			// process closing context
			return ErrContextClosed
		case payload = <-scannerDataCh:
			err := SendPayloadToWriter(payload, t.Conn)
			if err != nil {
				t.Close()

				return ErrPeerClosed
			}

			continue
		default:
		}

		// process err case after read all payload or Done context
		select {
		case err := <-errorCh:
			return err
		default:
		}
	}
}

func (t *Telnet) Receiver() error {
	scanner := bufio.NewScanner(t.Conn)
	scannerDataCh := make(chan []byte)

	go func() {
		// scanner.Scan() is lockable, moved to goroutine
		for scanner.Scan() {
			scannerDataCh <- scanner.Bytes()
		}
	}()

	var payload []byte
	for {
		select {
		case <-t.Ctx.Done():
			return ErrContextClosed
		case payload = <-scannerDataCh:
		}

		err := SendPayloadToWriter(payload, t.Writer)
		if err != nil {
			t.Close()

			return err
		}
	}
}

func SendPayloadToWriter(payload []byte, writer io.Writer) error {
	var buf bytes.Buffer
	buf.Write(payload)
	buf.WriteString("\n")

	_, err := writer.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
