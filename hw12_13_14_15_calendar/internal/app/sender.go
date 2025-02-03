package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/transport/ampq"
)

func RunSender(ctx context.Context, cfg *config.SenderConfig, log *slog.Logger) error {
	consumer := ampq.NewConsumer(log)
	if err := consumer.Setup(cfg.AmpqURI); err != nil {
		return fmt.Errorf("failed to setup consumer: %w", err)
	}

	deliveries, err := consumer.Consume()
	if err != nil {
		return fmt.Errorf("failed to consumer.Consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			if err := consumer.Shutdown(); err != nil {
				return fmt.Errorf("failed to shutdown consumer: %w", err)
			}
			return ctx.Err()
		case d := <-deliveries:
			var notification model.Notification

			err = json.Unmarshal(d.Body, &notification)
			if err != nil {
				log.Error("failed to json.Unmarshal", "body", string(d.Body), "error", err)
				continue
			}

			fmt.Printf("got notification %v", notification)

			d.Ack(false)
		}
	}
}
