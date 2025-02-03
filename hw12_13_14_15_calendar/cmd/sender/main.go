package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	cfg := config.NewSenderConfig()

	log, err := logger.SetupLogger(cfg.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("failed to setup logger: %+v", err))
	}

	go func() {
		if err = app.RunSender(ctx, cfg, log); err != nil {
			log.Error("critical service error", "error", err)
			stop()
			return
		}
	}()

	<-ctx.Done()

	log.Info("shutdown service ...")
}
