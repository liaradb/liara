package listener

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func LogGRPC() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		defer func() {
			d := time.Since(start)
			if err != nil {
				log.Printf("%v: %v\n\terror: %v\n", info.FullMethod, d, err)
			} else {
				log.Printf("%v: %v\n", info.FullMethod, d)
			}
		}()
		resp, err = handler(ctx, req)
		return
	}
}

func LogStreamGRPC() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		start := time.Now()
		defer func() {
			d := time.Since(start)
			if err != nil {
				log.Printf("%v: %v\n\terror: %v\n", info.FullMethod, d, err)
			} else {
				log.Printf("%v: %v\n", info.FullMethod, d)
			}
		}()
		err = handler(srv, stream)
		return
	}
}
