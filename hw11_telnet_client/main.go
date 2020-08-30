package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var ErrInterrupt = errors.New("Process was interrupted")

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: telnet <ADDRESS> <PORT>")
		os.Exit(1)
	}

	address, timeout, err := ParseArguments()
	if err != nil {
		fmt.Fprintf(os.Stderr, "...Incorrect arguments: %q\n", err)
		os.Exit(1)
	}

	c := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := c.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "...Connect error:", err)
		os.Exit(1)
	} else {
		fmt.Fprintln(os.Stderr, "...Connected to", address)
	}
	defer c.Close()
	errorCh := make(chan error)

	go func() {
		err := c.Sender()
		if err != nil {
			errorCh <- err
		}
	}()

	go func() {
		err := c.Receiver()
		if err != nil {
			errorCh <- err
		}
	}()

	// no wait signal goroutine
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT)

		select {
		case <-signalCh:
			errorCh <- ErrInterrupt
			c.Close()
		}
	}()

	err = <-errorCh
	fmt.Fprintf(os.Stderr, "...%s\n", err)
}

func ParseArguments() (string, time.Duration, error) {
	var timeoutFlag string

	flag.StringVar(&timeoutFlag, "timeout", "10s", "connect timeout in seconds")
	flag.Parse()
	timeout, err := time.ParseDuration(timeoutFlag)
	if err != nil {
		return "", timeout, err
	}

	var address bytes.Buffer
	address.WriteString(flag.Arg(0))
	address.WriteString(":")

	_, err = strconv.Atoi(flag.Arg(1))
	if err != nil {
		return "", timeout, errors.New("wrong port value")
	}
	address.WriteString(flag.Arg(1))

	return address.String(), timeout, nil
}
