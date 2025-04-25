package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_active_connections",
			Help: "The current number of active WebSocket connections",
		},
	)

	TotalConnections = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_connections_total",
			Help: "The total number of WebSocket connections established",
		},
	)

	WebsocketMessagesSent = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_messages_sent_total",
			Help: "The total number of WebSocket messages sent",
		},
	)

	ConnectionDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "websocket_connection_duration_seconds",
			Help:    "The duration of WebSocket connections in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	MessageSize = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "websocket_message_size_bytes",
			Help:    "Size of WebSocket messages in bytes",
			Buckets: []float64{64, 256, 1024, 4096, 16384, 65536},
		},
	)

	ApiRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_milliseconds",
			Help:    "API request latency in milliseconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	Errors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_errors_total",
			Help: "Total number of errors by type",
		}, []string{"type"},
	)
)
