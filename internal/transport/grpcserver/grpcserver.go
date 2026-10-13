package grpcserver

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Options struct {
	Host         string
	Port         int
	IdleTimeout  time.Duration
	Logger       *slog.Logger
	Interceptors []grpc.UnaryServerInterceptor
	Registers    []func(*grpc.Server)
}

type grpcServer struct {
	server *grpc.Server
	logger *slog.Logger
	opts   *Options
}

func NewGRPCServer(opts *Options) *grpcServer {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(opts.Interceptors...),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: opts.IdleTimeout,
		}),
	)

	for _, reg := range opts.Registers {
		reg(srv)
	}

	return &grpcServer{
		server: srv,
		logger: opts.Logger,
		opts:   opts,
	}
}

func (s *grpcServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.opts.Host, s.opts.Port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("grpcserver: failed to listen on %s: %w", addr, err)
	}

	s.logger.Info("starting gRPC server", "addr", addr)

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("grpcserver: serve error: %w", err)
	}

	s.logger.Info("gRPC server stopped gracefully")

	return nil
}

func (s *grpcServer) Stop() {
	s.server.GracefulStop()

	s.logger.Info("gRPC server stopped gracefully")
}
