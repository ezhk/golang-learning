package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrGoroutinesLimitNotPositive = errors.New("goroutines limit set to negative value")

type Task func() error

func Run(tasks []Task, concurrentGoroutinesLimit int, allowedErrorsAmount int) error {
	if concurrentGoroutinesLimit < 1 {
		return ErrGoroutinesLimitNotPositive
	}

	concurrentGoroutines := make(chan struct{}, concurrentGoroutinesLimit)
	defer close(concurrentGoroutines)

	for _, t := range tasks {
		if allowedErrorsAmount < 1 {
			return ErrErrorsLimitExceeded
		}

		concurrentGoroutines <- struct{}{}
		go func(t Task) {
			if err := t(); err != nil {
				allowedErrorsAmount--
			}
			<-concurrentGoroutines
		}(t)
	}

	// wait completed tasks
	for len(concurrentGoroutines) > 0 {
	}

	if allowedErrorsAmount < 1 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
