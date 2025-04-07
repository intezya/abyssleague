FROM golang:1.24.2-alpine3.21 AS builder

RUN apk add --no-cache git ca-certificates tzdata && \
    mkdir -p /build/services/lib/go /build/protos /build/services/websocket

WORKDIR /build

COPY services/lib/go/go.mod /build/services/lib/go/
COPY protos/go.mod /build/protos/
COPY services/websocket/go.mod services/websocket/go.sum* /build/services/websocket/

WORKDIR /build/services/websocket
RUN go mod download

COPY services/lib/go/ /build/services/lib/go/
COPY protos/ /build/protos/
COPY services/websocket/ /build/services/websocket/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-w -s -extldflags '-static'" \
    -o /build/websocket-service ./cmd

FROM alpine:3.19 AS debug

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app
COPY --from=builder /build/websocket-service .
ENV TZ=UTC

ENTRYPOINT ["/app/websocket-service"]

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/websocket-service /app/websocket-service

WORKDIR /app

ENTRYPOINT ["/app/websocket-service"]
