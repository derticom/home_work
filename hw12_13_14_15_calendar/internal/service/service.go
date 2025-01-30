package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"

	"github.com/google/uuid"
)

type Storage interface {
	Add(ctx context.Context, event model.Event) error
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, id model.EventUUID) error
	GetForDay(ctx context.Context, date time.Time) ([]model.Event, error)
	GetForWeek(ctx context.Context, date time.Time) ([]model.Event, error)
	GetForMonth(ctx context.Context, date time.Time) ([]model.Event, error)
}
type Service struct {
	storage Storage
	log     *slog.Logger
}

func New(storage Storage, log *slog.Logger) *Service {
	return &Service{
		storage: storage,
		log:     log,
	}
}

func (s *Service) Add(ctx context.Context, event model.Event) error {
	if uuid.UUID(event.ID) == uuid.Nil {
		event.ID = model.EventUUID(uuid.New())
	}
	return s.storage.Add(ctx, event)
}

func (s *Service) Update(ctx context.Context, event model.Event) error {
	return s.storage.Update(ctx, event)
}

func (s *Service) Delete(ctx context.Context, id model.EventUUID) error {
	return s.storage.Delete(ctx, id)
}

func (s *Service) GetForDay(ctx context.Context, date time.Time) ([]model.Event, error) {
	return s.storage.GetForDay(ctx, date)
}

func (s *Service) GetForWeek(ctx context.Context, date time.Time) ([]model.Event, error) {
	return s.storage.GetForWeek(ctx, date)
}

func (s *Service) GetForMonth(ctx context.Context, date time.Time) ([]model.Event, error) {
	return s.storage.GetForMonth(ctx, date)
}
