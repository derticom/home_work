// Package memorystorage - реализация in-memory хранилища.
package memorystorage

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
)

const week = 24 * 7 * time.Hour

type Events struct {
	eventsByUserID map[model.UserUUID][]model.Event // Ключ - UserID, значение - все его события.

	mu sync.RWMutex
}

func New() *Events {
	return &Events{
		eventsByUserID: make(map[model.UserUUID][]model.Event),
	}
}

// Add добавляет новое событие в хранилище.
func (e *Events) Add(_ context.Context, event model.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.eventsByUserID[event.UserID] = append(e.eventsByUserID[event.UserID], event)

	return nil
}

// Update обновляет уже имеющееся событие в хранилище у того же пользователя.
func (e *Events) Update(_ context.Context, id model.EventUUID, event model.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	userEvents, ok := e.eventsByUserID[event.UserID]
	if ok {
		for i, userEvent := range userEvents {
			if userEvent.ID == id {
				e.eventsByUserID[event.UserID] = slices.Delete(e.eventsByUserID[event.UserID], i, i+1)
				break
			}
		}
	}

	e.eventsByUserID[event.UserID] = append(e.eventsByUserID[event.UserID], event)

	return nil
}

// Delete удаляет событие из хранилища.
func (e *Events) Delete(_ context.Context, userID model.UserUUID, id model.EventUUID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	userEvents, ok := e.eventsByUserID[userID]
	if !ok || len(userEvents) == 0 {
		return nil
	}

	for i, userEvent := range userEvents {
		if userEvent.ID == id {
			e.eventsByUserID[userID] = slices.Delete(e.eventsByUserID[userID], i, i+1)
			return nil
		}
	}

	return nil
}

// GetForDay возвращает список событий по пользователю за указанный день.
func (e *Events) GetForDay(_ context.Context, userID model.UserUUID, date time.Time) ([]model.Event, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	var result []model.Event

	userEvents := e.eventsByUserID[userID]
	for _, event := range userEvents {
		if event.Date.Truncate(24 * time.Hour).Equal(date.Truncate(24 * time.Hour)) {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetForWeek возвращает список событий по пользователю за неделю (на входе - дата начала недели).
func (e *Events) GetForWeek(_ context.Context, userID model.UserUUID, date time.Time) ([]model.Event, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// date должен быть началом недели - понедельник.
	if date.Weekday() != time.Monday || date.Hour() != 0 || date.Minute() != 0 {
		return nil, errors.New("invalid date")
	}

	endPeriod := date.Add(week)

	var result []model.Event

	userEvents := e.eventsByUserID[userID]
	for _, event := range userEvents {
		if event.Date.After(date) && event.Date.Before(endPeriod) {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetForMonth возвращает список событий по пользователю за месяц.
func (e *Events) GetForMonth(_ context.Context, userID model.UserUUID, date time.Time) ([]model.Event, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	var result []model.Event

	userEvents := e.eventsByUserID[userID]
	for _, event := range userEvents {
		if event.Date.Year() == date.Year() &&
			event.Date.Month() == date.Month() {
			result = append(result, event)
		}
	}

	return result, nil
}
