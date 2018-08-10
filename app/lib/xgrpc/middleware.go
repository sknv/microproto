package xgrpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func WithLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	reply, err := handler(ctx, req)
	if err != nil {
		log.Printf("[ERROR] request %s failed with error: \"%s\" in %s", info.FullMethod, err, time.Since(start))
	} else {
		log.Printf("[INFO] request %s completed in %s", info.FullMethod, time.Since(start))
	}
	return reply, err
}
