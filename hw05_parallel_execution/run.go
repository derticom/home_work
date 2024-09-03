package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// CompletedTasksCount структура для подсчета выполненных задач
type CompletedTasksCount struct {
	count int
	mu    sync.Mutex
}

func (c *CompletedTasksCount) increase() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.count++
}

func (c *CompletedTasksCount) getCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.count
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		completedTasksCount CompletedTasksCount
		errCount            int // Счетчик ошибок.
		wg                  sync.WaitGroup
	)

	maxErrors := getMaxErrorsQuantity(m)

	done := make(chan struct{})
	tasksCh := make(chan Task)
	errCh := make(chan error)
	defer close(errCh)

	// Запуск воркеров.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(tasksCh, errCh, done, &completedTasksCount, &wg)
	}

	// Отправка задач воркерам.
	wg.Add(1)
	go func() {
		defer close(tasksCh)
		defer wg.Done()

		for _, task := range tasks {
			select {
			case <-done:
				break
			default:
				tasksCh <- task
			}
		}
	}()

	// Цикл с логикой управления остановкой функции.
	for {
		select {
		case <-errCh:
			errCount++
		default:
		}

		if maxErrors != 0 && errCount == maxErrors {
			break
		}

		if completedTasksCount.getCount() == len(tasks)-errCount {
			break
		}

		// Реализация остановки по граничному случаю из условия: "если в первых выполненных m задачах
		// (или вообще всех) происходят ошибки, то всего выполнится не более n+m задач."
		if errCount != 0 && completedTasksCount.getCount() == m+n {
			break
		}
	}

	close(done)
	wg.Wait()

	if maxErrors != 0 && errCount == maxErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(
	taskChanel chan Task,
	errChanel chan error,
	done chan struct{},
	count *CompletedTasksCount,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-done:
			return
		case task, ok := <-taskChanel:
			if !ok {
				continue
			}
			if err := task(); err != nil {
				select {
				case errChanel <- err:
				case <-done:
					return
				}
			} else {
				count.increase()
			}
		}
	}
}

func getMaxErrorsQuantity(m int) int {
	if m <= 0 {
		return 0
	}
	return m
}
