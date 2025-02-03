package grpc

import (
	"context"

	"google.golang.org/grpc"
)

func (s *Server) loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	s.serverLog.Info(info.FullMethod, req)

	return handler(ctx, req)
}
