package xgrpc

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHealthServer(grpcServer *grpc.Server) {
	healthv1.RegisterHealthServer(grpcServer, health.NewServer())
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type HealthServer struct {
	healthClient healthv1.HealthClient
}

func NewHealthServer(grpcConn *grpc.ClientConn) *HealthServer {
	return &HealthServer{healthClient: healthv1.NewHealthClient(grpcConn)}
}

func (s *HealthServer) Check(w http.ResponseWriter, _ *http.Request) {
	check, err := s.healthClient.Check(context.Background(), &healthv1.HealthCheckRequest{})
	if err != nil || check.GetStatus() != healthv1.HealthCheckResponse_SERVING {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("service unavailable"))
		return
	}
	w.Write([]byte("ok"))
}
