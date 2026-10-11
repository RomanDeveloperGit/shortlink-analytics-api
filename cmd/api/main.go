package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/config"
	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/infra/db"
	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/logger"
	linkRepository "github.com/RomanDeveloperGit/shortlink-analytics-api/internal/repository/link"
	linkService "github.com/RomanDeveloperGit/shortlink-analytics-api/internal/service/link"
	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/transport/grpcserver"
	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/transport/grpcserver/interceptor"
	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/transport/grpcserver/register"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	log := logger.MustSetup(&logger.Options{
		File:           os.Stdout,
		Env:            cfg.Env,
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
	})

	dbClient := db.MustConnect(&db.Options{
		ClientID:        cfg.ServiceName,
		ClientVersion:   cfg.ServiceVersion,
		Host:            cfg.ClickHouse.Host,
		Port:            cfg.ClickHouse.Port,
		User:            cfg.ClickHouse.User,
		Password:        cfg.ClickHouse.Password,
		DBName:          cfg.ClickHouse.DB,
		MaxIdleConns:    cfg.ClickHouse.MaxIdleConns,
		MaxOpenConns:    cfg.ClickHouse.MaxOpenConns,
		ConnMaxLifetime: cfg.ClickHouse.ConnMaxLifetime,
		DialTimeout:     cfg.ClickHouse.DialTimeout,
		PingTimeout:     cfg.ClickHouse.PingTimeout,
		Logger:          log,
	})

	linkRepo := linkRepository.NewRepository(dbClient.DB)
	linkSrv := linkService.NewService(linkRepo)

	registerAnalyticsService := register.NewRegisterAnalyticsService(linkSrv)

	grpcServer := grpcserver.NewGRPCServer(&grpcserver.Options{
		Host:        cfg.GRPCServer.Host,
		Port:        cfg.GRPCServer.Port,
		IdleTimeout: cfg.GRPCServer.IdleTimeout,
		Logger:      log,
		Registers:   []func(*grpc.Server){
			registerAnalyticsService,
		},
		Interceptors: []grpc.UnaryServerInterceptor{
			interceptor.Recovery(log),
			interceptor.TraceID(),
			interceptor.RequestID(),
			interceptor.Logging(log),
		},
	})

	time.Sleep(1 * time.Millisecond)

	go grpcServer.Start()

	log.Info("server started")

	<-ctx.Done()

	log.Info("shutting down server")

	// Останавливаем слушатель ОС сигналов, чтобы последующий второй CTRL+C мгновенно завершал работу.
	// Это работает, потому что снова произойдет INTERRUPT, на который никто не подписан(т.к. stop убирает эту подписку как раз)
	// А по умолчанию без подписки программа сразу же завершает своё исполнение
	// Сделано на случай, если нужно срочно остановить сервер, не дожидаясь полного graceful shutdown (часто в local-режиме)
	stop()

	grpcServer.Stop()
	dbClient.Close()

	log.Info("server stopped")
}
