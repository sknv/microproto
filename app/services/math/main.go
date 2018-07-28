package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/sknv/microproto/app/lib/xhttp"
	"github.com/sknv/microproto/app/services/math/cfg"
	"github.com/sknv/microproto/app/services/math/internal"
	"github.com/sknv/microproto/app/services/math/rpc"
)

const (
	shutdownTimeout = 60 * time.Second
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
	router.Use(middleware.RealIP, middleware.Logger)

	// handle twirp requests
	var srv internal.MathServer
	twirpHandler := rpc.NewMathServer(&srv, nil)
	router.Mount(rpc.MathPathPrefix, twirpHandler)

	// run the http server
	var healthCheck healthCheck
	router.Get("/healthz", healthCheck.healthz)
	xhttp.ListenAndServe(cfg.Addr, router, shutdownTimeout)
}
