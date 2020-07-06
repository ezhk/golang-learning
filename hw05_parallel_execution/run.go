package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded        = errors.New("errors limit exceeded")
	ErrGoroutinesLimitNonPositive = errors.New("goroutines limit contains non positive value")
	ErrErrorsAmountNonPositive    = errors.New("errors amount contains non positive value")
)

type Task func() error

func Run(tasks []Task, goroutinesLimit int, maxErrors int) error {
	if goroutinesLimit < 1 {
		return ErrGoroutinesLimitNonPositive
	}
	if maxErrors < 1 {
		return ErrErrorsAmountNonPositive
	}

	errorsLeft := int64(maxErrors)
	concurrentCh := make(chan struct{}, goroutinesLimit)
	defer close(concurrentCh)

	var wg sync.WaitGroup
	for _, task := range tasks {
		if atomic.LoadInt64(&errorsLeft) < 1 {
			break
		}

		// wait for next step: use channel as task queue
		concurrentCh <- struct{}{}

		wg.Add(1)
		go func(t Task) {
			defer func() {
				<-concurrentCh
				wg.Done()
			}()

			if err := t(); err != nil {
				atomic.AddInt64(&errorsLeft, -1)
			}
		}(task)
	}
	wg.Wait()

	if errorsLeft < 1 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
