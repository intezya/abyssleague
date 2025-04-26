package routes

import (
	"testing"
)

func TestIsInfoLogging(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Metrics path",
			path:     MetricsPath,
			expected: false,
		},
		{
			name:     "Ping path",
			path:     PingPath,
			expected: false,
		},
		{
			name:     "Websocket path prefix",
			path:     WebsocketPathPrefix,
			expected: true,
		},
		{
			name:     "Random path",
			path:     "/random",
			expected: true,
		},
		{
			name:     "Empty path",
			path:     "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := IsInfoLogging(tt.path)
			if result != tt.expected {
				t.Errorf("IsInfoLogging(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}
