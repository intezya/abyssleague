package routes

const (
	MetricsPath         = "/metrics"
	PingPath            = "/ping"
	WebsocketPathPrefix = "/websocket"
)

var NotInfoLogging = []string{PingPath, MetricsPath}

// IsInfoLogging if path is in NotInfoLogging return false and request must be logged as debug.
func IsInfoLogging(path string) bool {
	for _, v := range NotInfoLogging {
		if path == v {
			return false
		}
	}

	return true
}
