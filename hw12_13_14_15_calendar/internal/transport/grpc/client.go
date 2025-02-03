package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	log *slog.Logger

	conn *grpc.ClientConn
}

func NewClient(serverAddress string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to grpc.NewClient: %w", err)
	}

	return &Client{
		log: log,

		conn: conn,
	}, nil
}

func (c *Client) GetEvents(ctx context.Context) ([]*model.Event, error) {
	calendarClient := pb.NewCalendarClient(c.conn)

	req := pb.GetRequest{Date: &timestamppb.Timestamp{
		Seconds: time.Now().Truncate(24 * time.Hour).Unix(),
	}}

	response, err := calendarClient.GetForDay(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to calendarClient.GetForDay: %w", err)
	}

	return pbToModel(response.Events)
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		c.log.Error("failed to close gRPC client")
	}
}

func pbToModel(events []*pb.Event) ([]*model.Event, error) {
	pbEvents := make([]*model.Event, 0, len(events))

	for _, event := range events {
		id, err := uuid.Parse(event.Uuid)
		if err != nil {
			return nil, fmt.Errorf("failed to parse uuid: %w", err)
		}

		pbEvents = append(pbEvents, &model.Event{
			ID:           model.EventUUID(id),
			Header:       event.Header,
			Date:         time.Unix(event.Date.Seconds, int64(event.Date.Nanos)),
			Duration:     time.Duration(event.Duration),
			Description:  event.Description,
			NotifyBefore: time.Duration(event.NotifyBefore),
		})
	}

	return pbEvents, nil
}
