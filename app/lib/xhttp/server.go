package xhttp

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/sknv/microproto/app/lib/xos"
)

// ListenAndServe serves the handler at specified port
// and shuts down server gracefully in a shutdown timeout.
func ListenAndServe(addr string, handler http.Handler, shutdownTimeout time.Duration) {
	server := startServer(handler, addr)
	shutdownServerGracefully(server, shutdownTimeout)
}

func startServer(handler http.Handler, addr string) *http.Server {
	log.Print("[INFO] http server started on ", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Print("[ERROR] http server shutdown: ", err)
		}
	}()

	return server
}

func shutdownServerGracefully(server *http.Server, shutdownTimeout time.Duration) {
	// wait for interrupt signal to gracefully shutdown the server with a specified timeout
	xos.WaitForExit()

	log.Print("[INFO] shutting down the http server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	log.Print("[INFO] http server gracefully stopped")
}
