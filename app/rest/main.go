package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"

	"github.com/sknv/microproto/app/lib/xchi"
	"github.com/sknv/microproto/app/lib/xhttp"
	"github.com/sknv/microproto/app/rest/cfg"
	"github.com/sknv/microproto/app/rest/server"
)

const (
	concurrentRequestLimit = 1000
	serverShutdownTimeout  = 60 * time.Second
)

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type healthCheck struct{}

func (*healthCheck) healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

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

	// handle requests
	srv := server.NewRestServer(grpcConn)
	srv.Route(router)

	// start the http server
	var healthCheck healthCheck
	router.Get("/healthz", healthCheck.healthz)
	xhttp.ListenAndServe(cfg.Addr, router, serverShutdownTimeout)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("[FATAL] %s: %s", msg, err)
	}
}
