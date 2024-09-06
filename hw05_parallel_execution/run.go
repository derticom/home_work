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
		completedTasksCount int // Счетчик выполненных задач.
		errCount            int // Счетчик ошибок.
		wg                  sync.WaitGroup
	)

	maxErrors := getMaxErrorsQuantity(m)

	done := make(chan struct{})
	tasksCh := make(chan Task)
	completeStatusCh := make(chan bool)
	defer close(completeStatusCh)

	// Запуск воркеров.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(tasksCh, completeStatusCh, done, &wg)
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
		result := <-completeStatusCh
		if result {
			completedTasksCount++
		} else {
			errCount++
		}

		// Остановка при выполнении всех задач.
		if completedTasksCount == len(tasks)-errCount {
			break
		}

		// Остановка по достижению максимального количества ошибок.
		if maxErrors != 0 && errCount == maxErrors {
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
	completeStatusCh chan bool,
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
				case completeStatusCh <- false:
				case <-done:
					return
				}
			} else {
				select {
				case completeStatusCh <- true:
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
