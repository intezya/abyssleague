package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"os"
	"time"
)

type LokiConfig struct {
	url       string            // Example: "http://loki:3100/loki/api/v1/push"
	labels    map[string]string // labels is a tags for simple logs analysis
	batchSize int
	timeout   time.Duration
}

func NewLokiConfig(url string, labels map[string]string, batchSize int, timeout time.Duration) *LokiConfig {
	return &LokiConfig{url: url, labels: labels, batchSize: batchSize, timeout: timeout}
}

type lokiSink struct {
	config LokiConfig
	buffer []lokiEntry
	client *http.Client
}

type lokiEntry struct {
	Timestamp time.Time
	Line      string
}

type lokiPayload struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

func newLokiSink(config LokiConfig) zapcore.WriteSyncer {
	if config.batchSize <= 0 {
		config.batchSize = 100
	}
	if config.timeout <= 0 {
		config.timeout = 5 * time.Second
	}

	client := &http.Client{
		Timeout: config.timeout,
	}

	sink := &lokiSink{
		config: config,
		buffer: make([]lokiEntry, 0, config.batchSize),
		client: client,
	}

	go sink.periodicFlush()

	return sink
}

func (s *lokiSink) Write(p []byte) (n int, err error) {
	s.buffer = append(
		s.buffer, lokiEntry{
			Timestamp: time.Now(),
			Line:      string(p),
		},
	)

	if len(s.buffer) >= s.config.batchSize {
		go func() {
			err := s.flush()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "An error occurred while flushing logs to Loki: %v\n", err)
			}
		}()
	}

	return len(p), nil
}

func (s *lokiSink) Sync() error {
	return s.flush()
}

func (s *lokiSink) periodicFlush() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if len(s.buffer) > 0 {
			err := s.flush()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "An error occurred while flushing logs to Loki: %v\n", err)
			}
		}
	}
}

func (s *lokiSink) flush() error {
	if len(s.buffer) == 0 {
		return nil
	}

	entries := s.buffer
	s.buffer = make([]lokiEntry, 0, s.config.batchSize)

	values := make([][]string, 0, len(entries))
	for _, entry := range entries {
		timestampNano := fmt.Sprintf("%d", entry.Timestamp.UnixNano())
		values = append(values, []string{timestampNano, entry.Line})
	}

	payload := lokiPayload{
		Streams: []lokiStream{
			{
				Stream: s.config.labels,
				Values: values,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error marshaling Loki payload: %v\n", err)
		return err
	}

	req, err := http.NewRequest("POST", s.config.url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error creating Loki request: %v\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error sending logs to Loki: %v\n", err)
		return err
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	if resp.StatusCode >= 400 {
		_, _ = fmt.Fprintf(os.Stderr, "Loki API error, status code: %d\n", resp.StatusCode)
		return fmt.Errorf("loki API error: %d", resp.StatusCode)
	}

	return nil
}
