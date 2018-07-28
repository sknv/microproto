package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/sknv/microproto/app/lib/xchi"
	"github.com/sknv/microproto/app/lib/xhttp"
	"github.com/sknv/microproto/app/rest/cfg"
	"github.com/sknv/microproto/app/rest/server"
)

const (
	concurrentRequestLimit = 1000
	shutdownTimeout        = 60 * time.Second
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

	// config the http router
	router := chi.NewRouter()
	xchi.UseDefaultMiddleware(router)
	xchi.UseThrottle(router, concurrentRequestLimit)

	// handle requests
	srv := server.RestServer{Cfg: cfg}
	srv.Route(router)

	// run the http server
	var healthCheck healthCheck
	router.Get("/healthz", healthCheck.healthz)
	xhttp.ListenAndServe(cfg.Addr, router, shutdownTimeout)
}
