package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
	internalhttp "github.com/derticom/home_work/hw12_13_14_15_calendar/internal/server/http"
	srvс "github.com/derticom/home_work/hw12_13_14_15_calendar/internal/service"
	memorystorage "github.com/derticom/home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/derticom/home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface {
	Add(ctx context.Context, event model.Event) error
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, id model.EventUUID) error
	GetForDay(ctx context.Context, date time.Time) ([]model.Event, error)
	GetForWeek(ctx context.Context, date time.Time) ([]model.Event, error)
	GetForMonth(ctx context.Context, date time.Time) ([]model.Event, error)
}

func Run(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	var storage Storage
	if cfg.UseDataBaseStorage {
		dBStorage, err := sqlstorage.New(ctx, cfg.PostgresURL)
		if err != nil {
			return fmt.Errorf("failed to connect to postgres: %w", err)
		}

		err = dBStorage.Migrate("migrations")
		if err != nil {
			return fmt.Errorf("failed to dBStorage.Migrate: %w", err)
		}

		storage = dBStorage
	} else {
		storage = memorystorage.New()
	}

	service := srvс.New(storage, log)

	server := internalhttp.New(service, cfg.Server.Address, cfg.Server.Timeout, log)

	err := server.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to server.Run: %w", err)
	}

	return nil
}
