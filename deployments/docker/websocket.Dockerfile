FROM golang:1.24.2-alpine3.21 AS builder

RUN apk add --no-cache ca-certificates git tzdata  && \
    mkdir -p /build/services/lib/go /build/protos /build/services/websocket-api-gateway

WORKDIR /build

COPY protos/go.mod /build/protos/
COPY services/websocket/go.mod services/websocket/go.sum* /build/services/websocket-api-gateway/

WORKDIR /build/services/websocket-api-gateway
RUN go mod download

COPY protos/ /build/protos/
COPY services/websocket/ /build/services/websocket-api-gateway/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-w -s -extldflags '-static'" \
    -o /build/websocket-api-gateway ./cmd

FROM alpine:3.21 AS debug

RUN apk add --no-cache ca-certificates curl tzdata

WORKDIR /app

COPY --from=builder /build/websocket-api-gateway .

ENTRYPOINT ["/app/websocket-api-gateway"]

FROM alpine:3.21 AS final

RUN apk add --no-cache ca-certificates curl tzdata

COPY --from=builder /build/websocket-api-gateway /app/websocket-api-gateway

WORKDIR /app

ENTRYPOINT ["/app/websocket-api-gateway"]

USER nobody
