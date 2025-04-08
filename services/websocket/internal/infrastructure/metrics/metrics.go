package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActiveMainConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_active_main_connections",
			Help: "The current number of active main WebSocket connections",
		},
	)

	TotalMainConnections = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_main_connections_total",
			Help: "The total number of main WebSocket connections established",
		},
	)

	ActiveDraftConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_active_draft_connections",
			Help: "The current number of active draft WebSocket connections",
		},
	)

	TotalDraftConnections = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_draft_connections_total",
			Help: "The total number of draft WebSocket connections established",
		},
	)

	MainWebsocketMessagesSent = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "main_websocket_messages_sent_total",
			Help: "The total number of main WebSocket messages sent",
		},
	)

	DraftWebsocketMessagesSent = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "draft_websocket_messages_sent_total",
			Help: "The total number of draft WebSocket messages sent",
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
