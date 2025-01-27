// Package memorystorage - реализация in-memory хранилища.
package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
)

const week = 24 * 7 * time.Hour

type Events struct {
	events map[model.EventUUID]model.Event

	mu sync.RWMutex
}

func New() *Events {
	return &Events{
		events: make(map[model.EventUUID]model.Event),
	}
}

// Add добавляет новое событие в хранилище.
func (e *Events) Add(_ context.Context, event model.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.events[event.ID] = event

	return nil
}

// Update обновляет уже имеющееся событие в хранилище у того же пользователя.
func (e *Events) Update(_ context.Context, event model.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.events[event.ID] = event

	return nil
}

// Delete удаляет событие из хранилища.
func (e *Events) Delete(_ context.Context, id model.EventUUID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.events, id)

	return nil
}

// GetForDay возвращает список событий по пользователю за указанный день.
func (e *Events) GetForDay(_ context.Context, date time.Time) ([]model.Event, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	var result []model.Event

	for _, event := range e.events {
		if event.Date.Truncate(24 * time.Hour).Equal(date.Truncate(24 * time.Hour)) {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetForWeek возвращает список событий по пользователю за неделю (на входе - дата начала недели).
func (e *Events) GetForWeek(_ context.Context, date time.Time) ([]model.Event, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// date должен быть началом недели - понедельник.
	if date.Weekday() != time.Monday || date.Hour() != 0 || date.Minute() != 0 {
		return nil, errors.New("invalid date")
	}

	endPeriod := date.Add(week)

	var result []model.Event

	for _, event := range e.events {
		if event.Date.After(date) && event.Date.Before(endPeriod) {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetForMonth возвращает список событий по пользователю за месяц.
func (e *Events) GetForMonth(_ context.Context, date time.Time) ([]model.Event, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	var result []model.Event

	for _, event := range e.events {
		if event.Date.Year() == date.Year() &&
			event.Date.Month() == date.Month() {
			result = append(result, event)
		}
	}

	return result, nil
}
