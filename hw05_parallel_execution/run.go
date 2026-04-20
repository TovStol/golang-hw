package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, workers, maxErrors int) error {
	if maxErrors <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errCount int32
	taskCh := make(chan Task, len(tasks))
	var wg sync.WaitGroup

	go func() {
		for _, task := range tasks {
			taskCh <- task
		}
		close(taskCh)
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if atomic.LoadInt32(&errCount) >= int32(maxErrors) {
					return
				}
				if err := task(); err != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	if atomic.LoadInt32(&errCount) >= int32(maxErrors) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
