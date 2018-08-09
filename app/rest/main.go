package main

import (
	"log"
	"time"

	"github.com/go-chi/chi"
	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"

	"github.com/sknv/microproto/app/lib/xchi"
	"github.com/sknv/microproto/app/lib/xconsul"
	"github.com/sknv/microproto/app/lib/xhttp"
	"github.com/sknv/microproto/app/lib/xos"
	"github.com/sknv/microproto/app/rest/cfg"
	"github.com/sknv/microproto/app/rest/server"
)

const (
	concurrentRequestLimit = 1000
	serverShutdownTimeout  = 60 * time.Second

	serviceName         = "rest"
	healthCheckURL      = "/healthz"
	healthCheckInterval = "10s"
	healthCheckTimeout  = "1s"
)

func main() {
	cfg := cfg.Parse()

	// connect to grpc
	grpcConn, err := grpc.Dial(cfg.MathAddr, grpc.WithInsecure())
	xos.FailOnError(err, "failed to connect to grpc")
	defer grpcConn.Close()

	// config the http router
	router := chi.NewRouter()
	xchi.UseDefaultMiddleware(router)
	xchi.UseThrottle(router, concurrentRequestLimit)

	// handle requests
	rest := server.NewRestServer(grpcConn)
	rest.Route(router)

	// handle health check requests
	var health xhttp.HealthServer
	router.Get(healthCheckURL, health.Check)

	// start the http server and schedule a stop
	srv := xhttp.NewServer(cfg.Addr, router)
	srv.ListenAndServeAsync()
	defer srv.StopGracefully(serverShutdownTimeout)

	// register current service in consul and schedule a deregistration
	consulClient := registerConsulService(cfg)
	defer deregisterConsulService(consulClient)

	// wait for a program exit to stop the http server
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
		Name:     "rest api health check",
		HTTP:     "http://" + config.Addr + healthCheckURL,
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
