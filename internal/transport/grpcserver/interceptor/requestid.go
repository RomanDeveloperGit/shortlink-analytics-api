package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/pkg/observability"
)

func RequestID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		ctx, id := observability.ContextWithRequestID(ctx)

		_ = grpc.SetHeader(ctx, metadata.Pairs(observability.RequestIDKey, id))

		return handler(ctx, req)
	}
}
