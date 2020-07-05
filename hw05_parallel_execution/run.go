package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrGoroutinesLimitNonPositive = errors.New("goroutines limit contains non positive value")
var ErrErrorsAmountNonPositive = errors.New("errors amount contains non positive value")

type Task func() error

func Run(tasks []Task, goroutinesLimit int, allowedErrors int) error {
	if goroutinesLimit < 1 {
		return ErrGoroutinesLimitNonPositive
	}
	if allowedErrors < 1 {
		return ErrErrorsAmountNonPositive
	}

	errorsLeft := int64(allowedErrors)
	concurrentCh := make(chan struct{}, goroutinesLimit)

	for idx, task := range tasks {
		if errorsLeft < 1 {
			return ErrErrorsLimitExceeded
		}

		concurrentCh <- struct{}{}
		go func(t Task) {
			defer func() { <-concurrentCh }()

			if err := t(); err != nil {
				atomic.AddInt64(&errorsLeft, -1)
			}
		}(task)
	}

	// wait completed tasks and close chan
	for len(concurrentCh) > 0 {
	}
	close(concurrentCh)

	if errorsLeft < 1 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
