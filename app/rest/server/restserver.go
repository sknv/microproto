package server

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/sknv/microproto/app/lib/xgrpc"
	"github.com/sknv/microproto/app/lib/xhttp"
	math "github.com/sknv/microproto/app/services/math/rpc"
)

type RestServer struct {
	// consulClient *xconsul.Client
	mathClient math.MathClient
}

func NewRestServer(grpcConn *grpc.ClientConn) *RestServer {
	// consulClient, err := xconsul.NewClient(config.ConsulAddr)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to create a rest server")
	// }

	return &RestServer{
		// consulClient: consulClient,
		mathClient: math.NewMathClient(grpcConn),
	}
}

func (s *RestServer) Route(router chi.Router) {
	router.Route("/math", func(r chi.Router) {
		r.Get("/rect", s.Rect)
		r.Get("/circle", s.Circle)
	})
}

func (s *RestServer) Rect(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	width := parseFloat(w, queryParams.Get("w"))
	height := parseFloat(w, queryParams.Get("h"))
	args := math.RectArgs{
		Width:  width,
		Height: height,
	}

	reply, err := s.mathClient.Rect(context.Background(), &args)
	abortOnError(w, err)
	render.JSON(w, r, reply)
}

func (s *RestServer) Circle(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	radius := parseFloat(w, queryParams.Get("r"))
	args := math.CircleArgs{
		Radius: radius,
	}

	reply, err := s.mathClient.Circle(context.Background(), &args)
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

	gerr, _ := status.FromError(err)
	status := xgrpc.HTTPStatusFromCode(gerr.Code())
	if status != http.StatusInternalServerError {
		log.Print("[ERROR] abort on error: ", gerr.Err())
		http.Error(w, gerr.Message(), status)
		xhttp.AbortHandler()
	}
	panic(gerr)
}
