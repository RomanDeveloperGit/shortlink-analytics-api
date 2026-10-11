package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type database struct {
	DB     clickhouse.Conn
	logger *slog.Logger
}

type Options struct {
	ClientID        string
	ClientVersion   string
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	DialTimeout     time.Duration
	PingTimeout     time.Duration
	Logger          *slog.Logger
}

func MustConnect(opts *Options) *database {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", opts.Host, opts.Port)},
		Auth: clickhouse.Auth{
			Database: opts.DBName,
			Username: opts.User,
			Password: opts.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          opts.DialTimeout,
		MaxOpenConns:         opts.MaxOpenConns,
		MaxIdleConns:         opts.MaxIdleConns,
		ConnMaxLifetime:      opts.ConnMaxLifetime,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{
					Name:    opts.ClientID,
					Version: opts.ClientVersion,
				},
			},
		},
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.PingTimeout)
	defer cancel()

	err = conn.Ping(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to ping database: %w", err))
	}

	return &database{
		DB:     conn,
		logger: opts.Logger,
	}
}

func (d *database) Close() {
	err := d.DB.Close()
	if err != nil {
		d.logger.Error("failed to close database connection",
			slog.Any("error", err),
		)
	}
}
