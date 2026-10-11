package interceptor

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/pkg/observability"
)

const maxBodyLog = 1024

func getMD(md metadata.MD) string {
	var b strings.Builder

	for k, v := range md {
		b.WriteString(k)
		b.WriteString(":")
		b.WriteString(strings.Join(v, ","))
		b.WriteString(",")
	}

	return strings.TrimSuffix(b.String(), ",")
}

func marshalBody(m any) string {
	if m == nil {
		return ""
	}

	if pm, ok := m.(proto.Message); ok {
		b, err := protojson.Marshal(pm)
		if err == nil {
			return truncate(string(b))
		}
	}

	return truncate(slog.AnyValue(m).String())
}

func truncate(s string) string {
	if len(s) > maxBodyLog {
		return s[:maxBodyLog]
	}
	return s
}

func Logging(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		logger := observability.NewObservedLogger(ctx, log)

		var inMD metadata.MD
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			inMD = md
		}

		logger.Info("grpc request",
			slog.String("method", info.FullMethod),
			slog.String("headers", getMD(inMD)),
			slog.String("body", marshalBody(req)),
		)

		resp, err := handler(ctx, req)

		logger.Info("grpc response",
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
			slog.String("body", marshalBody(resp)),
			slog.Any("error", err),
		)

		return resp, err
	}
}
