package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"github.com/sknv/microproto/app/lib/xgrpc"
	"github.com/sknv/microproto/app/lib/xhttp"
	"github.com/sknv/microproto/app/lib/xos"
	"github.com/sknv/microproto/app/services/math/cfg"
	"github.com/sknv/microproto/app/services/math/internal"
	"github.com/sknv/microproto/app/services/math/rpc"
)

const (
	healthServerStartTimeout = 5 * time.Second
	serverShutdownTimeout    = 60 * time.Second
)

func main() {
	cfg := cfg.Parse()

	// start the grpc server and schedule a stop
	grpcSrv := startGrpcServerAsync(cfg)
	defer grpcSrv.StopGracefully(serverShutdownTimeout)

	time.Sleep(healthServerStartTimeout) // wait for grpc server to start

	// connect to grpc
	grpcConn, err := grpc.Dial(cfg.Addr, grpc.WithInsecure())
	failOnError(err, "failed to connect to grpc")
	defer grpcConn.Close()

	// start the health check server and schedule a stop
	healthSrv := startHealthServerAsync(cfg, grpcConn)
	defer healthSrv.StopGracefully(serverShutdownTimeout)

	// wait for a program exit to stop the health and grpc servers
	xos.WaitForExit()
}

func startGrpcServerAsync(config *cfg.Config) *xgrpc.Server {
	// listen on the specified address
	lis, err := net.Listen("tcp", config.Addr)
	failOnError(err, fmt.Sprintf("failed to listen on %s", config.Addr))

	// handle grpc requests
	srv := xgrpc.NewServer()
	xgrpc.RegisterHealthServer(srv.Server) // handle grpc health check requests
	rpc.RegisterMathServer(srv.Server, &internal.MathServer{})

	// start the grpc server
	srv.ServeAsync(lis)
	return srv
}

func startHealthServerAsync(config *cfg.Config, grpcConn *grpc.ClientConn) *xhttp.Server {
	// handle health check requests via http 1.1
	router := http.NewServeMux()
	health := xgrpc.NewHealthServer(grpcConn)
	router.HandleFunc("/healthz", health.Check)

	// start the http 1.1 health check server
	srv := xhttp.NewServer(config.HealthAddr, router)
	srv.ListenAndServeAsync()
	return srv
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("[FATAL] %s: %s", msg, err)
	}
}
