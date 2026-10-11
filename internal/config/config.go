package config

import (
	"time"

	"github.com/caarlos0/env/v11"

	internalEnv "github.com/RomanDeveloperGit/shortlink-analytics-api/internal/config/env"
)

type Config struct {
	Env            internalEnv.Env `env:"ENV,required,notEmpty"`
	ServiceName    string          `env:"SERVICE_NAME,required,notEmpty"`
	ServiceVersion string          `env:"SERVICE_VERSION,required,notEmpty"`
	GRPCServer     GRPCServer
	ClickHouse     ClickHouse
}

type GRPCServer struct {
	Host         string        `env:"GRPC_HOST,required,notEmpty"`
	Port         int           `env:"GRPC_PORT,required,notEmpty"`
	IdleTimeout  time.Duration `env:"GRPC_IDLE_TIMEOUT,required,notEmpty"`
}

type ClickHouse struct {
	Host            string        `env:"CLICKHOUSE_HOST,required,notEmpty"`
	Port            int           `env:"CLICKHOUSE_PORT,required,notEmpty"`
	User            string        `env:"CLICKHOUSE_USER,required,notEmpty"`
	Password        string        `env:"CLICKHOUSE_PASSWORD,required,notEmpty"`
	DB              string        `env:"CLICKHOUSE_DB,required,notEmpty"`
	MaxIdleConns    int           `env:"CLICKHOUSE_MAX_IDLE_CONNS,required,notEmpty"`
	MaxOpenConns    int           `env:"CLICKHOUSE_MAX_OPEN_CONNS,required,notEmpty"`
	ConnMaxLifetime time.Duration `env:"CLICKHOUSE_CONN_MAX_LIFETIME,required,notEmpty"`
	PingTimeout     time.Duration `env:"CLICKHOUSE_PING_TIMEOUT,required,notEmpty"`
	DialTimeout     time.Duration `env:"CLICKHOUSE_DIAL_TIMEOUT,required,notEmpty"`
}

func MustLoad() *Config {
	cfg := &Config{}

	// Специально через Must, чтобы вылезла паника в случае некорректного конфига
	cfg = env.Must(cfg, env.Parse(cfg))

	if !internalEnv.IsEnv(string(cfg.Env)) {
		panic("failed to load config: env is not supported")
	}

	return cfg
}
