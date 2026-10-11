package interceptor

import (
	"context"

	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/pkg/observability"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		var traceID string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(observability.TraceIDKey); len(vals) > 0 {
				traceID = vals[0]
			}
		}

		ctx, id := observability.ContextWithTraceID(ctx, traceID)

		_ = grpc.SetHeader(ctx, metadata.Pairs(observability.TraceIDKey, id))

		return handler(ctx, req)
	}
}
