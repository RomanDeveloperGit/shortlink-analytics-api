module github.com/RomanDeveloperGit/shortlink-analytics-api

go 1.26.4

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.49.0
	github.com/RomanDeveloperGit/shortlink-protos v0.0.0
	github.com/caarlos0/env/v11 v11.4.1
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

// по-нормальному надо через Buf Registry или через реальный гит репозиторий
replace github.com/RomanDeveloperGit/shortlink-protos => ../shortlink-protos

require (
	github.com/ClickHouse/ch-go v0.74.0 // indirect
	github.com/andybalholm/brotli v1.2.6 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/klauspost/compress v1.20.0 // indirect
	github.com/paulmach/orb v0.13.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.30 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	go.opentelemetry.io/otel v1.46.0 // indirect
	go.opentelemetry.io/otel/trace v1.46.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
)
