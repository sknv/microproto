package xgrpc

import (
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHealthServer(grpcServer *grpc.Server) {
	healthv1.RegisterHealthServer(grpcServer, health.NewServer())
}
