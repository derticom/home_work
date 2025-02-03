package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/transport/ampq"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/transport/grpc"

	"github.com/google/uuid"
)

func RunScheduler(ctx context.Context, cfg *config.SchedulerConfig, log *slog.Logger) error {
	producer, err := ampq.NewProducer(cfg.AmpqURI, log)
	if err != nil {
		return fmt.Errorf("failed to ampq.New: %w", err)
	}

	client, err := grpc.NewClient(cfg.GRPCAddress, log)
	if err != nil {
		return fmt.Errorf("failed to grpc.NewClient: %w", err)
	}
	defer client.Close()

	ticker := time.NewTicker(cfg.ProcessPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err = process(ctx, client, producer, cfg.ProcessPeriod, log); err != nil {
				return fmt.Errorf("failed to process: %w", err)
			}
		}
	}
}

func process(
	ctx context.Context,
	client *grpc.Client,
	producer *ampq.Producer,
	period time.Duration,
	log *slog.Logger,
) error {
	log.Info("starting process")
	events, err := client.GetEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to client.GetEvents: %w", err)
	}

	now := time.Now()

	for _, event := range events {
		notifyTime := event.Date.Add(-event.NotifyBefore)

		if now.Add(-period).Before(notifyTime) && notifyTime.Before(now.Add(period)) {
			notification := &model.Notification{
				ID:     uuid.UUID(event.ID),
				Header: event.Header,
				Date:   event.Date,
			}

			bytes, err := json.Marshal(notification)
			if err != nil {
				return fmt.Errorf("failed to marshal notification: %w", err)
			}

			if err := producer.Publish(bytes); err != nil {
				return fmt.Errorf("failed to producer.Publish: %w", err)
			}
		}
	}
	log.Info("end process")
	return nil
}
