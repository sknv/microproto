package main

import (
	"time"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"

	"github.com/sknv/microproto/app/lib/xchi"
	"github.com/sknv/microproto/app/lib/xhttp"
	"github.com/sknv/microproto/app/lib/xos"
	"github.com/sknv/microproto/app/rest/cfg"
	"github.com/sknv/microproto/app/rest/server"
)

const (
	concurrentRequestLimit = 1000
	serverShutdownTimeout  = 60 * time.Second
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
	rest, err := server.NewRestServer(cfg, grpcConn)
	xos.FailOnError(err, "failed to start the rest server")
	rest.Route(router)

	// handle health check requests
	var health xhttp.HealthServer
	router.Get("/healthz", health.Check)

	// start the http server and schedule a stop
	srv := xhttp.NewServer(cfg.Addr, router)
	srv.ListenAndServeAsync()
	defer srv.StopGracefully(serverShutdownTimeout)

	// wait for a program exit to stop the http server
	xos.WaitForExit()
}
