package xgrpc

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
}

func NewServer() *Server {
	srv := grpc.NewServer()
	return &Server{Server: srv}
}

func (s *Server) ServeAsync(listener net.Listener) {
	log.Print("[INFO] starting a grpc server on ", listener.Addr())
	go func() {
		if err := s.Serve(listener); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Print("[ERROR] failed to serve a grpc server: ", err)
		}
	}()
}

func (s *Server) StopGracefully(shutdownTimeout time.Duration) {
	log.Print("[INFO] stopping the grpc server...")

	// try to stop the server gracefuly
	serverStoppedGracefuly := make(chan struct{})
	go func() {
		s.GracefulStop()
		serverStoppedGracefuly <- struct{}{}
	}()

	// wait for a graceful shutdown and then stop the server forcibly
	select {
	case <-serverStoppedGracefuly:
		log.Print("[INFO] grpc server gracefully stopped")
	case <-time.After(shutdownTimeout):
		s.Stop()
		log.Print("[WARN] grpc server forcibly stopped")
	}
}
