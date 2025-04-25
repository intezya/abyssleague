# Websocket Messaging API Gateway

Websocket microservice that provides real-time communication between users.

## Table of Contents
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Architecture](#architecture)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

## Features
- Real-time bidirectional communication
- JWT authentication
- Multiple hubs support
- Logging with Loki
- Monitoring with Prometheus and Grafana
- Easy-Scalable architecture

## Requirements
- Go 1.18+ (required for generics support and improved interface features used in the codebase)
- Docker and Docker Compose
- Access to gRPC services (if applicable)

## Installation

### Using Docker Compose

```yaml
services:
  websocket-api-gateway:
    container_name: websocket-api-gateway
    build:
      context: .
      dockerfile: ./deployments/docker/websocket.Dockerfile
    ports:
      - "127.0.0.1:8090:8090"   # For local development
      - "127.0.0.1:50051:50051" # For local development
      - "127.0.0.1:50052:50052" # For local development
    environment:
      - HTTP_PORT=${WEBSOCKET_HTTP_PORT:-8090}
      - JWT_SECRET=${JWT_SECRET}
      - JWT_ISSUER=${JWT_ISSUER}
      - ENV_TYPE=${ENV_TYPE:-dev}
      - LOKI_ENDPOINT_URL=http://loki:3100/loki/api/v1/push
      - LOKI_LABELS={"app":"websocket-api-gateway", "env":"${ENV_TYPE}"}
      - GRPC_SERVER_PORTS=${GRPC_SERVER_PORT_START_FROM}
      - WEBSOCKET_HUBS=${WEBSOCKET_HUBS}
    healthcheck:
      test: [ "CMD", "curl", "-f", "http://localhost:8090/ping" ]
      interval: 30s
      timeout: 15s
      retries: 3
      start_period: 15s
```

### Local Development

Create `.env` file with required variables (see [Configuration](#configuration)).

Run the service: (from root of service)
```bash 
go run cmd/main.go
```

## Usage

Before seeing the usage, you might want to understand the [architecture](#architecture).

### WebSocket Connection

Connect to the WebSocket endpoint:
```
ws://<host>:<port>/websocket/{hub_type}
```

Websocket endpoint also requires jwt token in `Authorization` header 

## Configuration

The service can be configured using environment variables:

| Variable               | Description                                | Default                                  |
|------------------------|--------------------------------------------|------------------------------------------|
| `HTTP_PORT`            | HTTP server port                           | 8090                                     |
| `JWT_SECRET`           | Secret key for JWT authentication          | Required                                 |
| `JWT_ISSUER`           | Issuer for JWT tokens                      | "issuer"                                 |
| `ENV_TYPE`             | Environment type (dev, prod)               | "dev"                                    |
| `WEBSOCKET_HUBS`       | Comma-separated list of available hubs     | "main"                                   |
| `GRPC_SERVER_PORTS`    | Comma-separated list of ports for each hub | 50051                                    |
| `LOKI_ENDPOINT_URL`    | URL for Loki logging                       | "http://localhost:3100/loki/api/v1/push" |
| `LOKI_LABELS`          | JSON-encoded map of labels for Loki        | {}                                       |

## Architecture

Microservice uses shared local protos, so you need to build it with that protos.

### Modules structure:
```
.
├───protos
│   ├───src
│   └───websocket -- generated pb files
└───services
    ├───...
    └───websocket
        ├───cmd
        ├───domain
        ├───infrastructure
        ├───adapters
        └───...
```

### Path configuration
If you want to work with this, you should know that go.mod uses local paths:

```
replace abyssproto => ./../../protos/
```

It also takes in websocket.Dockerfile which is in ./deployments/docker/

The service follows a clean architecture approach, with directories like domain, infrastructure, adapters, etc.

## API Documentation

### Endpoints

- `GET /ping` - Health check endpoint
- `GET /websocket/{hub_type}` - WebSocket connection endpoint

## Development

### Building from source

```bash
go build -o websocket-api cmd/main.go (from root of domainservice)
```

### Generating protobuf files

```bash
make generate-proto (from root of repository)
```

## Monitoring

The service can be monitored using:

- **Prometheus**: Metrics collection (`http://<host>:9090`)
- **Grafana**: Visualization and dashboards (`http://<host>:3000`)
- **Loki**: Log aggregation

### Available metrics
- Connection count
- Message throughput
- Response times
- Error rates

## Troubleshooting

### Common issues

1. **Connection refused**: Check if the service is running and ports are correctly exposed.
2. **Authentication failed**: Verify JWT token is valid and not expired.
3. **Log shipping issues**: Ensure Loki URL is correct and accessible.
