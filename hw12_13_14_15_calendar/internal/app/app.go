package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
	internalgrpc "github.com/derticom/home_work/hw12_13_14_15_calendar/internal/server/grpc"
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

	grpcServer := internalgrpc.New(service, cfg.GRPCAddress, log)

	httpServer := internalhttp.New(service, cfg.HTTPServer.Address, cfg.HTTPServer.Timeout, log)

	errCh := make(chan error)

	go func() {
		if err := grpcServer.Run(); err != nil {
			errCh <- fmt.Errorf("failed to grpcServer.Run: %w", err)
		}
	}()

	go func() {
		if err := httpServer.Run(ctx); err != nil {
			errCh <- fmt.Errorf("failed to httpServer.Run: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("critical service error: %w", err)
	case <-ctx.Done():
	}

	return nil
}
