package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/pb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// При локальном тестировании calendar заменить на localhost.
const (
	restURL = "http://calendar:8085/"
	grpcURL = "calendar:8090"
)

func TestRestApi(t *testing.T) {
	testEvent := `
{
  "id": "ca8b788c-a0d0-471e-8264-d6e6b7d8d28a",
  "header": "Test Event",
  "date": "2025-01-28T12:00:00Z",
  "duration": 3600000000000,
  "description": "Test Event Description",
  "notify_before": 7200000000000
}
`

	// test "hello" handler
	resp, err := http.Get(restURL + "hello")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// test "add" handler
	resp, err = http.Post(restURL+"add", "Content-Type", strings.NewReader(testEvent))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// test "get_for_day" handler
	resp, err = http.Get(restURL + "get_for_day" + "?date=2025-01-28")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var events []model.Event
	err = json.Unmarshal(body, &events)
	require.NoError(t, err)
	assert.Equal(t, "ca8b788c-a0d0-471e-8264-d6e6b7d8d28a", uuid.UUID(events[0].ID).String())
	assert.Equal(t, "Test Event Description", events[0].Description)

	// test "delete" handler
	deleteRequest := `"ca8b788c-a0d0-471e-8264-d6e6b7d8d28a"`
	req, err := http.NewRequest(http.MethodDelete, restURL+"delete", strings.NewReader(deleteRequest))
	require.NoError(t, err)

	client := &http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "event deleted successfully")

	// Проверка отсутствия событий.
	resp, err = http.Get(restURL + "get_for_day" + "?date=2025-01-28")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	var eventsAfterDelete []model.Event
	err = json.Unmarshal(body, &eventsAfterDelete)
	require.NoError(t, err)
	assert.Empty(t, eventsAfterDelete)
}

func TestGrpcApi(t *testing.T) {
	// тестовые данные.
	todayEventUUID := uuid.New()
	today := time.Now().Add(5 * time.Hour)
	weekEventUUID := uuid.New()
	week := time.Now().Add(24 * 2 * time.Hour)
	monthEventUUID := uuid.New()
	month := time.Now().Add(24 * 10 * time.Hour)

	testEventToday := pb.Event{
		Uuid:         todayEventUUID.String(),
		Header:       "Test Event 1",
		Date:         timestamppb.New(today),
		Duration:     int64(100 * time.Second),
		Description:  "Test Event 1 Description",
		NotifyBefore: int64(100 * time.Second),
	}

	testEventWeek := pb.Event{
		Uuid:         weekEventUUID.String(),
		Header:       "Test Event 2",
		Date:         timestamppb.New(week),
		Duration:     int64(100 * time.Second),
		Description:  "Test Event 2 Description",
		NotifyBefore: int64(100 * time.Second),
	}

	testEventMonth := pb.Event{
		Uuid:         monthEventUUID.String(),
		Header:       "Test Event 3",
		Date:         timestamppb.New(month),
		Duration:     int64(100 * time.Second),
		Description:  "Test Event 3 Description",
		NotifyBefore: int64(100 * time.Second),
	}

	testErrEvent := pb.Event{
		Uuid:         "invalid_uuid",
		Header:       "Test Event Err",
		Date:         timestamppb.New(time.Now()),
		Duration:     int64(100 * time.Second),
		Description:  "Test Event Err Description",
		NotifyBefore: int64(100 * time.Second),
	}

	// Создание подключения и клиента.
	conn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCalendarClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Тестирование добавления событий
	_, err = client.Add(ctx, &testEventToday)
	require.NoError(t, err)

	_, err = client.Add(ctx, &testEventWeek)
	require.NoError(t, err)

	_, err = client.Add(ctx, &testEventMonth)
	require.NoError(t, err)

	_, err = client.Add(ctx, &testErrEvent) // кейс с ошибкой - некорректный UUID.
	require.Error(t, err)
	st, _ := status.FromError(err)
	require.Equal(t, codes.InvalidArgument, st.Code())

	// Получение листинга событий на день/неделю/месяц.
	dayEvents, err := client.GetForDay(
		ctx, &pb.GetRequest{
			Date: timestamppb.New(today.Truncate(24 * time.Hour)),
		})
	require.NoError(t, err)
	assert.Equal(t, dayEvents.Events[0].Uuid, todayEventUUID.String())
	assert.Equal(t, dayEvents.Events[0].Header, testEventToday.Header)

	weekEvents, err := client.GetForWeek(
		ctx, &pb.GetRequest{
			Date: timestamppb.New(today.Truncate(24 * time.Hour)),
		})
	require.NoError(t, err)
	assert.Len(t, weekEvents.Events, 2)

	monthEvents, err := client.GetForMonth(
		ctx, &pb.GetRequest{
			Date: timestamppb.New(today.Truncate(24 * time.Hour)),
		})
	require.NoError(t, err)
	assert.Len(t, monthEvents.Events, 3)
}
