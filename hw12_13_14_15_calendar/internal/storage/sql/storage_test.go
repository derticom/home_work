//go:build integration

package sqlstorage_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
	sqlstorage "github.com/derticom/home_work/hw12_13_14_15_calendar/internal/storage/sql"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	dsn := "postgresql://postgres:password@localhost:5432/test?sslmode=disable"

	storage, err := sqlstorage.New(ctx, dsn)
	require.NoError(t, err)

	defer func(storage *sqlstorage.Storage, ctx context.Context) {
		err := storage.Close()
		require.NoError(t, err)
	}(storage, ctx)

	eventUUID := model.EventUUID(uuid.New())
	date := time.Date(2025, 1, 25, 12, 0, 0, 0, time.Local)
	testEvent := model.Event{
		ID:           eventUUID,
		Header:       "Some event header",
		Date:         date,
		Duration:     30 * time.Minute,
		Description:  "Some description of event",
		NotifyBefore: 3 * time.Hour,
	}

	testEventUpd := model.Event{
		ID:           eventUUID,
		Header:       "Some event header",
		Date:         date,
		Duration:     99 * time.Minute,
		Description:  "Some description of event",
		NotifyBefore: 1 * time.Hour,
	}

	// Проверка добавления.
	err = storage.Add(ctx, testEvent)
	require.NoError(t, err)
	gotEvents, err := storage.GetForDay(ctx, date)
	require.NoError(t, err)
	require.Equal(t, testEvent, gotEvents[0])

	// Проверка обновления.
	err = storage.Update(ctx, testEventUpd)
	require.NoError(t, err)
	gotEventsUpdated, err := storage.GetForDay(ctx, date)
	require.Equal(t, testEventUpd, gotEventsUpdated[0])

	// Проверка удаления.
	err = storage.Delete(ctx, eventUUID)
	require.NoError(t, err)
	got, err := storage.GetForDay(ctx, date)
	require.NoError(t, err)
	require.Empty(t, got)
}

func Test_getMonthRange(t *testing.T) {
	janStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
	febStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.Local)

	type args struct {
		date time.Time
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name: "january 2025",
			args: args{
				date: janStart,
			},
			wantStart: janStart,
			wantEnd:   time.Date(2025, 1, 31, 23, 59, 59, 999999999, time.Local),
		},
		{
			name: "february 2025",
			args: args{
				date: febStart,
			},
			wantStart: febStart,
			wantEnd:   time.Date(2025, 2, 28, 23, 59, 59, 999999999, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := sqlstorage.GetMonthRange(tt.args.date)
			if !reflect.DeepEqual(gotStart, tt.wantStart) {
				t.Errorf("getMonthRange() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if !reflect.DeepEqual(gotEnd, tt.wantEnd) {
				t.Errorf("getMonthRange() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}
