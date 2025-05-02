package errorz

import "time"

// ErrorResponse is the JSON structure returned to clients
type ErrorResponse struct {
	Message     string                 `json:"message"`
	Detail      string                 `json:"detail,omitempty"`
	Code        int                    `json:"code"`
	Path        string                 `json:"path"`
	Timestamp   time.Time              `json:"timestamp"`
	ErrorID     string                 `json:"error_id"`
	Type        string                 `json:"type"`
	Validations []string               `json:"validations,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
