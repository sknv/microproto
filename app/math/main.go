package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	consul "github.com/hashicorp/consul/api"

	"github.com/sknv/microproto/app/lib/xconsul"
	"github.com/sknv/microproto/app/lib/xhttp"
	"github.com/sknv/microproto/app/lib/xos"
	"github.com/sknv/microproto/app/math/cfg"
	"github.com/sknv/microproto/app/math/rpc"
	"github.com/sknv/microproto/app/math/server"
)

const (
	serverShutdownTimeout = 60 * time.Second

	serviceName         = "math"
	healthCheckURL      = "/healthz"
	healthCheckInterval = "10s"
	healthCheckTimeout  = "1s"
)

func main() {
	cfg := cfg.Parse()

	// config the http router
	router := chi.NewRouter()
	router.Use(middleware.RealIP, middleware.Logger)

	// handle requests
	var math server.MathServer
	twirpHandler := rpc.NewMathServer(&math, nil)
	router.Mount(rpc.MathPathPrefix, twirpHandler)

	// handle health check requests
	var health xhttp.HealthServer
	router.Get("/healthz", health.Check)

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

	tags := []string{fmt.Sprintf("urlprefix-/%s strip=/%s", serviceName, serviceName)} // for fabio load balancer
	healthCheck := &consul.AgentServiceCheck{
		Name:     "math service health check",
		HTTP:     "http://" + config.Addr + healthCheckURL,
		Interval: healthCheckInterval,
		Timeout:  healthCheckTimeout,
	}
	if err = consulClient.RegisterCurrentService(
		config.Addr, serviceName, tags, consul.AgentServiceChecks{healthCheck},
	); err != nil {
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
