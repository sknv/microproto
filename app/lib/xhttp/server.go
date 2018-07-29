package xhttp

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/sknv/microproto/app/lib/xos"
)

// ListenAndServe serves the handler on the specified address
// and stops the server gracefully with the specified timeout.
func ListenAndServe(addr string, handler http.Handler, shutdownTimeout time.Duration) {
	server := startServer(handler, addr)
	stopServerGracefully(server, shutdownTimeout)
}

func startServer(handler http.Handler, addr string) *http.Server {
	log.Print("[INFO] starting http server on ", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Print("[ERROR] http server stopped: ", err)
		}
	}()
	return server
}

func stopServerGracefully(server *http.Server, shutdownTimeout time.Duration) {
	// wait for a program exit to stop the server gracefully with the specified timeout
	xos.WaitForExit()

	log.Print("[INFO] stopping the http server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("[FATAL] failed to stop the http server gracefully: ", err)
	}
	log.Print("[INFO] http server gracefully stopped")
}
