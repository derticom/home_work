// Package sqlstorage - реализация хранилища в БД Postgres.
package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib" //revive:disable:blank-imports // import for side effect need here.
	"github.com/pressly/goose/v3"
)

const week = 24 * 7 * time.Hour

type Storage struct {
	db *sql.DB
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to sql.Open: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to db.PingContext: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to Close: %w", err)
	}

	return nil
}

func (s *Storage) Migrate(migrate string) (err error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to goose.SetDialect: %w", err)
	}

	if err := goose.Up(s.db, migrate); err != nil {
		return fmt.Errorf("failed to goose.Up: %w", err)
	}

	return nil
}

func (s *Storage) Add(ctx context.Context, event model.Event) error {
	query := `
		INSERT INTO events (id, header, date, duration, description, notify_before)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.ExecContext(ctx, query,
		uuid.UUID(event.ID),
		event.Header,
		event.Date,
		event.Duration.Milliseconds(),
		event.Description,
		event.NotifyBefore.Milliseconds(),
	)
	if err != nil {
		return fmt.Errorf("failed to db.ExecContext: %w", err)
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, event model.Event) error {
	query := `
	UPDATE events SET 
		header = $1,
		date = $2, 
		duration = $3,
		description = $4,
		notify_before = $5
    WHERE id = $6`

	_, err := s.db.ExecContext(ctx, query,
		event.Header,
		event.Date,
		event.Duration.Milliseconds(),
		event.Description,
		event.NotifyBefore.Milliseconds(),
		uuid.UUID(event.ID),
	)
	if err != nil {
		return fmt.Errorf("failed to db.ExecContext: %w", err)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id model.EventUUID) error {
	query := `DELETE FROM events WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, uuid.UUID(id))
	if err != nil {
		return fmt.Errorf("failed to db.ExecContext: %w", err)
	}

	return nil
}

func (s *Storage) GetForDay(ctx context.Context, date time.Time) ([]model.Event, error) {
	query := `
		SELECT id, header, date, duration, description, notify_before	
		FROM events
		WHERE DATE(date) = $1`

	rows, err := s.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to db.QueryContext: %w", err)
	}
	defer rows.Close()

	events, err := scan(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows: %w", err)
	}

	return events, nil
}

func (s *Storage) GetForWeek(ctx context.Context, date time.Time) ([]model.Event, error) {
	query := `
		SELECT id, header, date, duration, description, notify_before	
		FROM events
		WHERE date >= $1 AND date < $2`

	rows, err := s.db.QueryContext(ctx, query, date, date.Add(week))
	if err != nil {
		return nil, fmt.Errorf("failed to db.QueryContext: %w", err)
	}
	defer rows.Close()

	events, err := scan(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows: %w", err)
	}

	return events, nil
}

func (s *Storage) GetForMonth(ctx context.Context, date time.Time) ([]model.Event, error) {
	start, end := GetMonthRange(date)

	query := `
		SELECT id, header, date, duration, description, notify_before	
		FROM events
		WHERE date >= $1 AND date < $2`

	rows, err := s.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to db.QueryContext: %w", err)
	}
	defer rows.Close()

	events, err := scan(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows: %w", err)
	}

	return events, nil
}

func scan(rows *sql.Rows) ([]model.Event, error) {
	var events []model.Event

	for rows.Next() {
		var id uuid.UUID
		var header string
		var dateInfo time.Time
		var duration int64
		var description string
		var notifyBefore int64

		err := rows.Scan(
			&id,
			&header,
			&dateInfo,
			&duration,
			&description,
			&notifyBefore,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to rows.Scan: %w", err)
		}

		events = append(events, model.Event{
			ID:           model.EventUUID(id),
			Header:       header,
			Date:         dateInfo,
			Duration:     time.Duration(duration) * time.Millisecond,
			Description:  description,
			NotifyBefore: time.Duration(notifyBefore) * time.Millisecond,
		})
	}

	return events, nil
}

// GetMonthRange - возвращает период - месяц. На входе - начало месяца.
func GetMonthRange(date time.Time) (start, end time.Time) {
	startOfNextMonth := date.AddDate(0, 1, 0)

	return date, startOfNextMonth.Add(-time.Nanosecond)
}
