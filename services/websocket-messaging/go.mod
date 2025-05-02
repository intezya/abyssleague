module github.com/intezya/abyssleague/services/websocket-messaging

go 1.24.2

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/intezya/abyssleague/proto v0.0.0
	github.com/intezya/pkglib/configloader v0.1.2
	github.com/intezya/pkglib/itertools v0.1.1
	github.com/intezya/pkglib/logger v0.1.2
	github.com/prometheus/client_golang v1.22.0
	google.golang.org/grpc v1.71.1
	google.golang.org/protobuf v1.36.6
)

replace github.com/intezya/abyssleague/proto => ./../../protos/

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)
