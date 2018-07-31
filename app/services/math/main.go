package main

import (
	"fmt"
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
	serverShutdownTimeout = 60 * time.Second
	// serviceName           = "math"
)

func main() {
	cfg := cfg.Parse()

	// start the grpc server and schedule a stop
	grpcSrv := startGrpcServerAsync(cfg)
	defer grpcSrv.StopGracefully(serverShutdownTimeout)

	// connect to grpc
	grpcConn, err := grpc.Dial(cfg.Addr, grpc.WithInsecure())
	xos.FailOnError(err, "failed to connect to grpc")
	defer grpcConn.Close()

	// start the health check server and schedule a stop
	healthSrv := startHealthServerAsync(cfg, grpcConn)
	defer healthSrv.StopGracefully(serverShutdownTimeout)

	// register current service in consul and schedule a deregistration
	//
	// consulClient := registerConsulService(cfg)
	// defer deregisterConsulService(consulClient)

	// wait for a program exit to stop the health and grpc servers
	xos.WaitForExit()
}

func startGrpcServerAsync(config *cfg.Config) *xgrpc.Server {
	// listen on the specified address
	lis, err := net.Listen("tcp", config.Addr)
	xos.FailOnError(err, fmt.Sprintf("failed to listen on %s", config.Addr))

	// handle grpc requests
	srv := xgrpc.NewServer()
	rpc.RegisterMathServer(srv.Server, &internal.MathServer{})
	xgrpc.RegisterHealthServer(srv.Server) // handle grpc health check requests

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

// consul section
//
// func registerConsulService(config *cfg.Config) *xconsul.Client {
// 	consulClient, err := xconsul.NewClient(config.ConsulAddr)
// 	if err != nil {
// 		log.Print("[ERROR] failed to connect to consul: ", err)
// 		return nil
// 	}

// 	if err = consulClient.RegisterService(config.Addr, serviceName); err != nil {
// 		log.Print("[ERROR] failed to register current service: ", err)
// 		return nil
// 	}
// 	return consulClient
// }

// func deregisterConsulService(consulClient *xconsul.Client) {
// 	if consulClient == nil {
// 		return
// 	}

// 	if err := consulClient.DeregisterService(); err != nil {
// 		log.Print("[ERROR] failed to deregister current service: ", err)
// 	}
// }
