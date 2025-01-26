package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Server struct {
	address string
	timeout time.Duration

	log       *slog.Logger
	serverLog *slog.Logger
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func New(
	address string,
	timeout time.Duration,
	log *slog.Logger,
) *Server {
	file, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Error("failed to os.OpenFile", "error", err)
		file = os.Stderr
	}

	serverLog := slog.New(slog.NewJSONHandler(file, nil))

	return &Server{
		address:   address,
		timeout:   timeout,
		log:       log,
		serverLog: serverLog,
	}
}

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.handleMain)

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

	s.log.Info("server listening on " + s.address)
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

func (s *Server) handleMain(w http.ResponseWriter, r *http.Request) {
	s.log.Info("start processing", "request", r.URL.Path)

	_, err := w.Write([]byte("hello world"))
	if err != nil {
		s.log.Error("failed to write response", "error", err)
		_, err := fmt.Fprintf(w, "failed to write response")
		if err != nil {
			s.log.Error("failed to write response", "error", err)
		}
		return
	}

	s.log.Info("successfully finished processing", "request", r.URL.Path)
}
