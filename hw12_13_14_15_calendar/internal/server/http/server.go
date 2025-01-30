package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
)

type Server struct {
	service Service
	address string
	timeout time.Duration

	log       *slog.Logger
	serverLog *slog.Logger
}

type Service interface {
	Add(ctx context.Context, event model.Event) error
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, id model.EventUUID) error
	GetForDay(ctx context.Context, date time.Time) ([]model.Event, error)
	GetForWeek(ctx context.Context, date time.Time) ([]model.Event, error)
	GetForMonth(ctx context.Context, date time.Time) ([]model.Event, error)
}

func New(
	service Service,
	address string,
	timeout time.Duration,
	log *slog.Logger,
) *Server {
	file, err := os.OpenFile("http_server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Error("failed to os.OpenFile", "error", err)
		file = os.Stderr
	}

	serverLog := slog.New(slog.NewJSONHandler(file, nil))

	return &Server{
		service:   service,
		address:   address,
		timeout:   timeout,
		log:       log,
		serverLog: serverLog,
	}
}

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", s.main)
	mux.HandleFunc("POST /add", s.add)
	mux.HandleFunc("PUT /update", s.update)
	mux.HandleFunc("DELETE /delete", s.delete)
	mux.HandleFunc("GET /get_for_day", s.getForDay)
	mux.HandleFunc("GET /get_for_week", s.getForWeek)
	mux.HandleFunc("GET /get_for_month", s.getForMonth)

	loggedMux := s.loggingMiddleware(mux)

	server := http.Server{
		Addr:        s.address,
		Handler:     loggedMux,
		ReadTimeout: s.timeout,
	}

	go func() {
		<-ctx.Done()
		err := server.Shutdown(ctx)
		if err != nil {
			s.log.Error("failed to server.Shutdown", "error", err)
		}
	}()

	s.log.Info("http server listening on " + s.address)
	err := server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("shutdown server")
			return nil
		}
		return fmt.Errorf("failed to server.ListenAndServe: %w", err)
	}

	return nil
}
