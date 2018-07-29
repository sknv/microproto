package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/sknv/microproto/app/lib/xgrpc"
	"github.com/sknv/microproto/app/services/math/cfg"
	"github.com/sknv/microproto/app/services/math/internal"
	"github.com/sknv/microproto/app/services/math/rpc"
)

const (
	serverShutdownTimeout = 60 * time.Second
)

// TODO: health check

func main() {
	cfg := cfg.Parse()

	// listen on the specified address
	lis, err := net.Listen("tcp", cfg.Addr)
	failOnError(err, fmt.Sprintf("failed to listen on %s", cfg.Addr))

	// handle grpc requests
	srv := grpc.NewServer()
	rpc.RegisterMathServer(srv, &internal.MathServer{})

	// start the grpc server
	xgrpc.Serve(srv, lis, serverShutdownTimeout)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("[FATAL] %s: %s", msg, err)
	}
}
