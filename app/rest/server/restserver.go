package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/twitchtv/twirp"

	"github.com/sknv/microproto/app/lib/xhttp"
	"github.com/sknv/microproto/app/lib/xtwirp"
	math "github.com/sknv/microproto/app/math/rpc"
	"github.com/sknv/microproto/app/rest/cfg"
)

type RestServer struct {
	mathClient math.Math
}

func NewRestServer(config *cfg.Config) *RestServer {
	return &RestServer{mathClient: math.NewMathProtobufClient(config.MathProxyAddr, &http.Client{})}
}

func (s *RestServer) Route(router chi.Router) {
	router.Route("/math", func(r chi.Router) {
		r.Get("/rect", s.rect)
		r.Get("/circle", s.circle)
	})
}

func (s *RestServer) rect(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	width := parseFloat(w, queryParams.Get("w"))
	height := parseFloat(w, queryParams.Get("h"))
	args := math.RectArgs{
		Width:  width,
		Height: height,
	}

	// set the reply timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := s.mathClient.Rect(ctx, &args)
	abortOnError(w, err)
	render.JSON(w, r, reply)
}

func (s *RestServer) circle(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	radius := parseFloat(w, queryParams.Get("r"))
	args := math.CircleArgs{
		Radius: radius,
	}

	// set the reply timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := s.mathClient.Circle(ctx, &args)
	abortOnError(w, err)
	render.JSON(w, r, reply)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func parseFloat(w http.ResponseWriter, s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Print("[ERROR] parse float: ", err)
		http.Error(w, "argument must be a float number", http.StatusBadRequest)
		xhttp.AbortHandler()
	}
	return val
}

func abortOnError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	log.Print("[ERROR] abort on error: ", err)

	cause := errors.Cause(err)
	status, _ := xtwirp.FromError(cause)
	httpCode := twirp.ServerHTTPStatusFromErrorCode(status.Code())
	if httpCode != http.StatusInternalServerError {
		http.Error(w, status.Msg(), httpCode)
		xhttp.AbortHandler()
	}
	xhttp.AbortHandlerWithInternalError(w)
}
