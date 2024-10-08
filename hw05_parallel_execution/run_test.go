package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	//nolint:depguard // Применение 'require' необходимо для тестирования.
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("without errors", func(t *testing.T) {
		taskList := []Task{
			func() error { return nil },
			func() error { return nil },
			func() error { return nil },
		}

		err := Run(taskList, 2, 1)
		require.NoError(t, err)
	})

	t.Run("with errors but below the limit", func(t *testing.T) {
		taskList := []Task{
			func() error { return nil },
			func() error { return errors.New("error 1") },
			func() error { return nil },
			func() error { return errors.New("error 2") },
			func() error { return nil },
		}

		err := Run(taskList, 3, 3)
		require.NoError(t, err)
	})

	t.Run("errors exceed limit", func(t *testing.T) {
		taskList := []Task{
			func() error { return nil },
			func() error { return errors.New("error 1") },
			func() error { return errors.New("error 2") },
			func() error { return nil },
			func() error { return errors.New("error 3") },
		}

		err := Run(taskList, 2, 2)
		require.Error(t, err)
		require.Equal(t, ErrErrorsLimitExceeded, err)
	})

	t.Run("all tasks have errors", func(t *testing.T) {
		taskList := []Task{
			func() error { return errors.New("error 1") },
			func() error { return errors.New("error 2") },
			func() error { return errors.New("error 3") },
		}

		err := Run(taskList, 2, 1)
		require.Error(t, err)
		require.Equal(t, ErrErrorsLimitExceeded, err)
	})
}

func Test_getMaxErrorsQuantity(t *testing.T) {
	type args struct {
		m int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"five",
			args{5},
			5,
		},
		{
			"zero",
			args{0},
			0,
		},
		{
			"negative",
			args{-5},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMaxErrorsQuantity(tt.args.m); got != tt.want {
				t.Errorf("getMaxErrorsQuantity() = %v, want %v", got, tt.want)
			}
		})
	}
}
