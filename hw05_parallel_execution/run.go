package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		completedTasksCount int  // Счетчик выполненных задач.
		errCount            int  // Счетчик ошибок.
		boundaryCase        bool // Флаг наличия ошибок при выполнении m задач.
		wg                  sync.WaitGroup
	)

	maxErrors := getMaxErrorsQuantity(m)

	done := make(chan struct{})
	tasksCh := make(chan Task)
	completedTasksCh := make(chan struct{})
	defer close(completedTasksCh)
	errCh := make(chan error)
	defer close(errCh)

	// Запуск воркеров.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(tasksCh, completedTasksCh, errCh, done, &wg)
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
				select {
				case tasksCh <- task:
				case <-done:
					return
				}
			}
		}
	}()

	// Цикл с логикой управления остановкой функции.
	for {
		select {
		case <-completedTasksCh:
			completedTasksCount++
		case <-errCh:
			errCount++
		}

		// Остановка при достижении макс. кол-ва ошибок.
		if maxErrors != 0 && errCount == maxErrors {
			break
		}

		// Остановка при выполнении всех задач.
		if completedTasksCount == len(tasks)-errCount {
			break
		}

		// Остановка по граничному случаю из условия: "если в первых выполненных m задачах
		// (или вообще всех) происходят ошибки, то всего выполнится не более n+m задач."
		if maxErrors != 0 && errCount != 0 && completedTasksCount == m {
			boundaryCase = true
		}
		if boundaryCase && completedTasksCount+errCount == m+n {
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
	completedTasksCh chan struct{},
	errChanel chan error,
	done chan struct{},
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
				select {
				case completedTasksCh <- struct{}{}:
				case <-done:
					return
				}
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
