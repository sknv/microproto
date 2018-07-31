package main

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

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

	// config the http router
	router := chi.NewRouter()
	router.Use(middleware.RealIP, middleware.Logger)

	// handle requests
	var math internal.MathServer
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
	//
	// consulClient := registerConsulService(cfg)
	// defer deregisterConsulService(consulClient)

	// wait for a program exit to stop the http server
	xos.WaitForExit()
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
