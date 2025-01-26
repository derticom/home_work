package memorystorage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
	memorystorage "github.com/derticom/home_work/hw12_13_14_15_calendar/internal/storage/memory"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	eventsStorage := memorystorage.New()

	userUUID := model.UserUUID(uuid.New())
	eventUUID := model.EventUUID(uuid.New())
	date := time.Date(2025, 1, 25, 12, 0, 0, 0, time.UTC)
	dateNextMonth := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println(dateNextMonth)
	event := model.Event{
		ID:           eventUUID,
		Header:       "Some event header",
		Date:         date,
		Duration:     3 * time.Hour,
		Description:  "Some description of event",
		UserID:       userUUID,
		NotifyBefore: 3 * time.Hour,
	}

	eventUpdated := model.Event{
		ID:           eventUUID,
		Header:       "Some event header",
		Date:         date.Add(24 * 7 * 2 * time.Hour),
		Duration:     4 * time.Hour,
		Description:  "Some description of event",
		UserID:       userUUID,
		NotifyBefore: 1 * time.Hour,
	}

	// Добавление и проверка получения события на заданный день.
	err := eventsStorage.Add(ctx, event)
	require.NoError(t, err)
	gotForDay, err := eventsStorage.GetForDay(ctx, userUUID, time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	assert.Equal(t, []model.Event{event}, gotForDay)

	// Проверка получения на неделю.
	gotForWeek, err := eventsStorage.GetForWeek(ctx, userUUID, time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	assert.Equal(t, []model.Event{event}, gotForWeek)

	// Изменение события - перенос на 2 недели вперед.
	err = eventsStorage.Update(ctx, event.ID, eventUpdated)
	require.NoError(t, err)
	gotForMonth, err := eventsStorage.GetForMonth(ctx, userUUID, dateNextMonth)
	require.NoError(t, err)
	require.Len(t, gotForMonth, 1)
	assert.Equal(t, eventUpdated, gotForMonth[0])

	// Удаление события.
	err = eventsStorage.Delete(ctx, userUUID, eventUUID)
	require.NoError(t, err)
	gotForMonth, err = eventsStorage.GetForMonth(ctx, userUUID, date)
	require.NoError(t, err)
	assert.Empty(t, gotForMonth)
}
