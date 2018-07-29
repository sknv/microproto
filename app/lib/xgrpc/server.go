package xgrpc

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/sknv/microproto/app/lib/xos"
)

// Serve listens and serves the grpc server and stops it gracefully with the specified timeout.
func Serve(server *grpc.Server, listener net.Listener, shutdownTimeout time.Duration) {
	startServer(server, listener)
	stopServerGracefully(server, shutdownTimeout)
}

func startServer(server *grpc.Server, listener net.Listener) {
	log.Print("[INFO] starting grpc server on ", listener.Addr())
	go func() {
		if err := server.Serve(listener); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Print("[ERROR] grpc server stopped: ", err)
		}
	}()
}

func stopServerGracefully(server *grpc.Server, shutdownTimeout time.Duration) {
	// wait for a program exit to stop the server gracefully with the specified timeout
	xos.WaitForExit()

	log.Print("[INFO] stopping the grpc server...")

	// wait for a graceful shutdown and then stop the server forcibly
	shutdownTimer := time.NewTimer(shutdownTimeout)
	go func() {
		<-shutdownTimer.C
		server.Stop()
		log.Print("[WARN] grpc server forcibly stopped")
	}()

	// try to stop the server gracefuly
	server.GracefulStop()
	serverStoppedGracefuly := shutdownTimer.Stop()
	if serverStoppedGracefuly {
		log.Print("[INFO] grpc server gracefully stopped")
	}
}
