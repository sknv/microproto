package server

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/twitchtv/twirp"

	"github.com/sknv/microproto/app/lib/xhttp"
	"github.com/sknv/microproto/app/rest/cfg"
	math "github.com/sknv/microproto/app/services/math/rpc"
)

type RestServer struct {
	Cfg *cfg.Config
}

func (s *RestServer) Route(router chi.Router) {
	router.Get("/math/rect", s.Rect)
	router.Get("/math/circle", s.Circle)
}

func (s *RestServer) Rect(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	width := parseFloat(w, queryParams.Get("w"))
	height := parseFloat(w, queryParams.Get("h"))
	args := math.RectArgs{
		Width:  width,
		Height: height,
	}

	mathClient := math.NewMathProtobufClient(s.Cfg.MathURL, &http.Client{})
	reply, err := mathClient.Rect(context.Background(), &args)
	failOnError(w, err)
	render.JSON(w, r, reply)
}

func (s *RestServer) Circle(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	radius := parseFloat(w, queryParams.Get("r"))
	args := math.CircleArgs{
		Radius: radius,
	}

	mathClient := math.NewMathProtobufClient(s.Cfg.MathURL, &http.Client{})
	reply, err := mathClient.Circle(context.Background(), &args)
	failOnError(w, err)
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

func failOnError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	twerr := err.(twirp.Error)
	errStatus := twirp.ServerHTTPStatusFromErrorCode(twerr.Code())
	if errStatus != http.StatusInternalServerError {
		log.Print("[ERROR] ", twerr)
		http.Error(w, twerr.Error(), errStatus)
		xhttp.AbortHandler()
	}
	panic(twerr)
}
