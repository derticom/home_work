package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	Add(ctx context.Context, event model.Event) error
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, id model.EventUUID) error
	GetForDay(ctx context.Context, date time.Time) ([]model.Event, error)
	GetForWeek(ctx context.Context, date time.Time) ([]model.Event, error)
	GetForMonth(ctx context.Context, date time.Time) ([]model.Event, error)
}

type Server struct {
	pb.UnimplementedCalendarServer

	service Service
	address string

	log       *slog.Logger
	serverLog *slog.Logger
}

func New(service Service, address string, log *slog.Logger) *Server {
	file, err := os.OpenFile("grpc_server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Error("failed to os.OpenFile", "error", err)
		file = os.Stderr
	}

	serverLog := slog.New(slog.NewJSONHandler(file, nil))

	return &Server{
		service:   service,
		address:   address,
		log:       log,
		serverLog: serverLog,
	}
}

func (s *Server) Run() error {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(s.loggingInterceptor),
	)
	pb.RegisterCalendarServer(grpcServer, s)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to net.Listen: %w", err)
	}

	s.log.Info("grpc server listening on " + s.address)

	return grpcServer.Serve(listener)
}

func (s *Server) Add(ctx context.Context, event *pb.Event) (*emptypb.Empty, error) {
	id, err := uuid.Parse(event.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse uuid: %v", err)
	}

	date := time.Unix(event.GetDate().Seconds, 0)

	modelEvent := model.Event{
		ID:           model.EventUUID(id),
		Header:       event.GetHeader(),
		Date:         date,
		Duration:     time.Duration(event.Duration),
		Description:  event.GetDescription(),
		NotifyBefore: time.Duration(event.NotifyBefore),
	}

	if err = s.service.Add(ctx, modelEvent); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add event: %v", err)
	}

	return nil, nil
}

func (s *Server) Update(ctx context.Context, event *pb.Event) (*emptypb.Empty, error) {
	id, err := uuid.Parse(event.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse uuid: %v", err)
	}

	date := time.Unix(event.GetDate().Seconds, 0)

	modelEvent := model.Event{
		ID:           model.EventUUID(id),
		Header:       event.GetHeader(),
		Date:         date,
		Duration:     time.Duration(event.Duration),
		Description:  event.GetDescription(),
		NotifyBefore: time.Duration(event.NotifyBefore),
	}

	if err = s.service.Update(ctx, modelEvent); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update event: %v", err)
	}

	return nil, nil
}

func (s *Server) Delete(ctx context.Context, request *pb.DelRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(request.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse uuid: %v", err)
	}

	err = s.service.Delete(ctx, model.EventUUID(id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete event: %v", err)
	}

	return nil, nil
}

func (s *Server) GetForDay(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	if request.GetDate() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "date is required")
	}

	date := time.Unix(request.GetDate().Seconds, 0)

	events, err := s.service.GetForDay(ctx, date.Truncate(24*time.Hour))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}

	return &pb.GetResponse{
		Events: modelToPb(events),
	}, nil
}

func (s *Server) GetForWeek(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	if request.GetDate() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "date is required")
	}

	date := time.Unix(request.GetDate().Seconds, 0)

	events, err := s.service.GetForWeek(ctx, date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}

	return &pb.GetResponse{
		Events: modelToPb(events),
	}, nil
}

func (s *Server) GetForMonth(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	if request.GetDate() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "date is required")
	}

	date := time.Unix(request.GetDate().Seconds, 0)

	events, err := s.service.GetForMonth(ctx, date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}

	return &pb.GetResponse{
		Events: modelToPb(events),
	}, nil
}

func modelToPb(events []model.Event) []*pb.Event {
	pbEvents := make([]*pb.Event, 0, len(events))

	for _, event := range events {
		pbEvents = append(pbEvents, &pb.Event{
			Uuid:         uuid.UUID(event.ID).String(),
			Header:       event.Header,
			Date:         timestamppb.New(event.Date),
			Duration:     int64(event.Duration.Seconds()),
			Description:  event.Description,
			NotifyBefore: int64(event.NotifyBefore.Seconds()),
		})
	}

	return pbEvents
}
