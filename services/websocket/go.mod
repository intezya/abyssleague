module websocket

go 1.24.2

require go.uber.org/zap v1.27.0 // indirect

require abysslib v0.0.0

require (
	abyssproto v0.0.0
	github.com/fasthttp/websocket v1.5.12
	google.golang.org/grpc v1.71.1
	google.golang.org/protobuf v1.36.6
)

replace abysslib => ./../lib/go

replace abyssproto => ./../../protos/

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/savsgio/gotils v0.0.0-20240704082632-aef3928b8a38 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.58.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
)
