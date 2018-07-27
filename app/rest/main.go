package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/sknv/microproto/app/lib/xchi"
	"github.com/sknv/microproto/app/lib/xhttp"
	"github.com/sknv/microproto/app/rest/cfg"
)

const (
	concurrentRequestLimit = 1000
	shutdownTimeout        = 60 * time.Second
)

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type healthcheck struct{}

func (*healthcheck) healthz(w http.ResponseWriter, _ *http.Request) {
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

	// route the server
	// var srv server.Server
	// srv.Route(router)

	// run the http server
	var healthcheck healthcheck
	router.Get("/healthz", healthcheck.healthz)
	xhttp.ListenAndServe(cfg.Addr, router, shutdownTimeout)
}
