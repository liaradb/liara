package listener

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func LogGRPC() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		defer func() {
			d := time.Since(start)
			if err != nil {
				slog.Error("request",
					"method", info.FullMethod,
					"time", d,
					"error", err)
			} else {
				slog.Info("request",
					"method", info.FullMethod,
					"time", d)
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
				slog.Error("request",
					"method", info.FullMethod,
					"time", d,
					"error", err)
			} else {
				slog.Info("request",
					"method", info.FullMethod,
					"time", d)
			}
		}()
		err = handler(srv, stream)
		return
	}
}
