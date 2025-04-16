FROM golang:1.24.2-alpine3.21 AS builder

RUN apk add --no-cache ca-certificates git tzdata  && \
    mkdir -p /build/services/lib/go /build/protos /build/services/abysscore

WORKDIR /build

COPY protos/go.mod /build/protos/
COPY services/abysscore/go.mod services/abysscore/go.sum* /build/services/abysscore/

WORKDIR /build/services/abysscore
RUN go mod download

COPY protos/ /build/protos/
COPY services/abysscore/ /build/services/abysscore/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-w -s -extldflags '-static'" \
    -o /build/abysscore ./cmd

FROM alpine:3.21 AS debug

RUN apk add --no-cache ca-certificates curl tzdata

WORKDIR /app

COPY --from=builder /build/services/abysscore .

ENTRYPOINT ["/app/abysscore"]

FROM alpine:3.21 AS final

RUN apk add --no-cache ca-certificates curl tzdata

COPY --from=builder /build/abysscore /app/abysscore

WORKDIR /app

ENTRYPOINT ["/app/abysscore"]

USER nobody
