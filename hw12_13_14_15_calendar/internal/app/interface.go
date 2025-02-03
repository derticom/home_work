package app

import (
	"context"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
)

type Storage interface {
	Add(ctx context.Context, event model.Event) error
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, id model.EventUUID) error
	GetForDay(ctx context.Context, date time.Time) ([]model.Event, error)
	GetForWeek(ctx context.Context, date time.Time) ([]model.Event, error)
	GetForMonth(ctx context.Context, date time.Time) ([]model.Event, error)
}
