package main

import (
	"log"
	"net"
	"time"

	consul "github.com/hashicorp/consul/api"

	"github.com/sknv/microproto/app/lib/xconsul"
	"github.com/sknv/microproto/app/lib/xgrpc"
	"github.com/sknv/microproto/app/lib/xos"
	"github.com/sknv/microproto/app/math/cfg"
	"github.com/sknv/microproto/app/math/rpc"
	"github.com/sknv/microproto/app/math/server"
)

const (
	serverShutdownTimeout = 60 * time.Second

	serviceName         = "math"
	healthCheckInterval = "10s"
	healthCheckTimeout  = "1s"
)

func main() {
	cfg := cfg.Parse()

	// listen on the specified address
	lis, err := net.Listen("tcp", cfg.Addr)
	xos.FailOnError(err, "failed to listen on "+cfg.Addr)

	// handle grpc requests
	srv := xgrpc.NewServer()
	rpc.RegisterMathServer(srv.Server, &server.MathServer{})
	xgrpc.RegisterHealthServer(srv.Server) // handle grpc health check requests

	// start the grpc server and schedule a stop
	srv.ServeAsync(lis)
	defer srv.StopGracefully(serverShutdownTimeout)

	// register current service in consul and schedule a deregistration
	consulClient := registerConsulService(cfg)
	defer deregisterConsulService(consulClient)

	// wait for a program exit to stop the health and grpc servers
	xos.WaitForExit()
}

// ----------------------------------------------------------------------------
// consul section
// ----------------------------------------------------------------------------

func registerConsulService(config *cfg.Config) *xconsul.Client {
	consulClient, err := xconsul.NewClient(config.ConsulAddr)
	if err != nil {
		log.Print("[ERROR] failed to connect to consul: ", err)
		return nil
	}

	healthCheck := &consul.AgentServiceCheck{
		Name:     "math service health check",
		GRPC:     config.Addr,
		Interval: healthCheckInterval,
		Timeout:  healthCheckTimeout,
	}
	if err = consulClient.RegisterCurrentService(config.Addr, serviceName, healthCheck); err != nil {
		log.Print("[ERROR] failed to register current service: ", err)
		return nil
	}
	return consulClient
}

func deregisterConsulService(consulClient *xconsul.Client) {
	if consulClient == nil {
		return
	}

	if err := consulClient.DeregisterCurrentService(); err != nil {
		log.Print("[ERROR] failed to deregister current service: ", err)
	}
}
