package main

import (
	"log"
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
	failOnError(err, "failed to connect to grpc")
	defer grpcConn.Close()

	// config the http router
	router := chi.NewRouter()
	xchi.UseDefaultMiddleware(router)
	xchi.UseThrottle(router, concurrentRequestLimit)

	// handle health check requests
	var health xhttp.HealthServer
	router.Get("/healthz", health.Check)

	// handle requests
	rest := server.NewRestServer(grpcConn)
	rest.Route(router)

	// start the http server and schedule a stop
	srv := xhttp.NewServer(cfg.Addr, router)
	srv.ListenAndServeAsync()
	defer srv.StopGracefully(serverShutdownTimeout)

	// wait for a program exit to stop the http server
	xos.WaitForExit()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("[FATAL] %s: %s", msg, err)
	}
}
