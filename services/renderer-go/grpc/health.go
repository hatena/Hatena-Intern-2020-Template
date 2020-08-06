package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// gRPC Health Checking Protocol の実装
// https://github.com/grpc/grpc/blob/master/doc/health-checking.md

// Check はサービスの稼働状況を返す
func (s *Server) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	if in.Service == "" || in.Service == "renderer.Renderer" {
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
	} else {
		return nil, status.Error(codes.NotFound, "service not found")
	}
}

// func (s *Server) Watch(req *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
// 	return nil
// }
