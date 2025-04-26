package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestMetricsInitialization(t *testing.T) {
	t.Parallel()
	// Test that all metrics are properly initialized
	// Test simple metrics
	simpleMetrics := []struct {
		name   string
		metric interface {
			Desc() *prometheus.Desc
		}
	}{
		{"ActiveConnections", ActiveConnections},
		{"TotalConnections", TotalConnections},
		{"WebsocketMessagesSent", WebsocketMessagesSent},
		{"ConnectionDuration", ConnectionDuration},
		{"MessageSize", MessageSize},
	}

	for _, m := range simpleMetrics {
		if m.metric.Desc() == nil {
			t.Errorf("Metric %s is not properly initialized", m.name)
		}
	}

	// Test vector metrics by using them with labels
	// This will panic if they're not properly initialized
	ApiRequestDuration.WithLabelValues("GET", "/test")
	Errors.WithLabelValues("test_error")
}

func TestMetricsUsage(t *testing.T) {
	t.Parallel()
	// Test that we can use the metrics without errors
	// This is a simple smoke test to ensure the metrics are usable
	// Test incrementing counters
	TotalConnections.Inc()
	WebsocketMessagesSent.Inc()

	// Test setting gauge
	ActiveConnections.Set(42)

	// Test observing histograms
	ConnectionDuration.Observe(0.5)
	MessageSize.Observe(1024)

	// Test using histogram vector
	ApiRequestDuration.WithLabelValues("GET", "/api/test").Observe(0.1)

	// Test using counter vector
	Errors.WithLabelValues("test_error").Inc()
}
