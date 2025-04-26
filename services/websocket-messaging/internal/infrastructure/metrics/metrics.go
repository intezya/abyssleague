package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "websocket",
			Subsystem:   "connections",
			Name:        "websocket_active_connections",
			Help:        "The current number of active WebSocket connections",
			ConstLabels: nil,
		},
	)

	TotalConnections = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "websocket",
			Subsystem:   "connections",
			Name:        "websocket_connections_total",
			Help:        "The total number of WebSocket connections established",
			ConstLabels: nil,
		},
	)

	WebsocketMessagesSent = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "websocket",
			Subsystem:   "messages",
			Name:        "websocket_messages_sent_total",
			Help:        "The total number of WebSocket messages sent",
			ConstLabels: nil,
		},
	)

	ConnectionDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace:                       "websocket",
			Subsystem:                       "connections",
			Name:                            "websocket_connection_duration_seconds",
			Help:                            "The duration of WebSocket connections in seconds",
			Buckets:                         prometheus.DefBuckets,
			ConstLabels:                     nil,
			NativeHistogramBucketFactor:     0,
			NativeHistogramZeroThreshold:    0,
			NativeHistogramMaxBucketNumber:  0,
			NativeHistogramMinResetDuration: 0,
			NativeHistogramMaxZeroThreshold: 0,
			NativeHistogramMaxExemplars:     0,
			NativeHistogramExemplarTTL:      0,
		},
	)

	MessageSize = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace:                       "websocket",
			Subsystem:                       "messages",
			Name:                            "websocket_message_size_bytes",
			Help:                            "Size of WebSocket messages in bytes",
			Buckets:                         []float64{64, 256, 1024, 4096, 16384, 65536},
			ConstLabels:                     nil,
			NativeHistogramBucketFactor:     0,
			NativeHistogramZeroThreshold:    0,
			NativeHistogramMaxBucketNumber:  0,
			NativeHistogramMinResetDuration: 0,
			NativeHistogramMaxZeroThreshold: 0,
			NativeHistogramMaxExemplars:     0,
			NativeHistogramExemplarTTL:      0,
		},
	)

	ApiRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:                       "api",
			Subsystem:                       "http",
			Name:                            "api_request_duration_seconds",
			Help:                            "API request latency in seconds",
			Buckets:                         prometheus.DefBuckets,
			ConstLabels:                     nil,
			NativeHistogramBucketFactor:     0,
			NativeHistogramZeroThreshold:    0,
			NativeHistogramMaxBucketNumber:  0,
			NativeHistogramMinResetDuration: 0,
			NativeHistogramMaxZeroThreshold: 0,
			NativeHistogramMaxExemplars:     0,
			NativeHistogramExemplarTTL:      0,
		},
		[]string{"method", "endpoint"},
	)

	Errors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "websocket",
			Subsystem:   "errors",
			Name:        "websocket_errors_total",
			Help:        "Total number of errors by type",
			ConstLabels: nil,
		}, []string{"type"},
	)
)
