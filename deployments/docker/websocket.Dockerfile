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
    -o /build/websocket-api-gateway ./cmd

FROM alpine:3.21 AS debug

RUN apk add --no-cache tzdata ca-certificates curl

WORKDIR /app
COPY --from=builder /build/websocket-api-gateway .
ENV TZ=UTC

ENTRYPOINT ["/app/websocket-api-gateway"]

FROM alpine:3.21 AS final

RUN apk add --no-cache tzdata ca-certificates curl

COPY --from=builder /build/websocket-api-gateway /app/websocket-api-gateway

WORKDIR /app

ENTRYPOINT ["/app/websocket-api-gateway"]
