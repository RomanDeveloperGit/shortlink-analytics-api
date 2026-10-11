package interceptor

import (
	"context"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("recovered from panic",
					slog.Any("error", rec),
					slog.String("method", info.FullMethod),
					slog.String("stack", string(debug.Stack())),
				)

				err = status.Errorf(
					codes.Internal,
					"internal server error",
				)
			}
		}()

		return handler(ctx, req)
	}
}
