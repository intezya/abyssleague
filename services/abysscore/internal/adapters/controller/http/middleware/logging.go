package middleware

import (
	"abysscore/internal/adapters/config"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/pkglib/logger"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type LoggingMiddleware struct {
	config               *config.Config
	metrics              *prometheusMetrics
	resourceSampler      *resourceSampler
	slowRequestThreshold time.Duration
}

type prometheusMetrics struct {
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
}

type resourceSampler struct {
	requestCounter int
	samplingRate   int
	mutex          sync.Mutex
}

var textContentTypes = []string{
	"text/",
	"application/json",
	"application/xml",
	"application/javascript",
	"application/xhtml+xml",
}

func NewLoggingMiddleware(config *config.Config) *LoggingMiddleware {
	metrics := initPrometheusMetrics()

	slowRequestThreshold := 500 * time.Millisecond
	if config.SlowRequestThresholdMs > 0 {
		slowRequestThreshold = time.Duration(config.SlowRequestThresholdMs) * time.Millisecond
	}

	return &LoggingMiddleware{
		config:               config,
		metrics:              metrics,
		resourceSampler:      newResourceSampler(100),
		slowRequestThreshold: slowRequestThreshold,
	}
}

func newResourceSampler(samplingRate int) *resourceSampler {
	return &resourceSampler{
		requestCounter: 0,
		samplingRate:   samplingRate,
	}
}

func initPrometheusMetrics() *prometheusMetrics {
	metrics := &prometheusMetrics{
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		requestSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "Size of HTTP requests in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path"},
		),
		responseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "Size of HTTP responses in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path", "status"},
		),
	}

	prometheus.MustRegister(metrics.requestDuration, metrics.requestSize, metrics.responseSize)

	return metrics
}

func (l *LoggingMiddleware) Handle() fiber.Handler {
	tracer := otel.Tracer("fiber-middleware")

	return func(c *fiber.Ctx) error {
		requestID := c.Locals(l.config.FiberRequestIDConfig.ContextKey)

		ctx, span := tracer.Start(c.UserContext(), fmt.Sprintf("%s %s", c.Method(), c.Path()))
		defer span.End()

		c.SetUserContext(ctx)

		setSpanAttributes(span, c, requestID)

		logger.Log.Debugf("Request ID from context: '%v'", requestID)

		start := time.Now()
		err := c.Next()
		requestDuration := time.Since(start)

		statusCode := c.Response().StatusCode()

		responseBodyStr := extractResponseBody(c)

		log := createLogEntry(c, requestID, statusCode, requestDuration, responseBodyStr, err)

		l.updateMetrics(c, statusCode, requestDuration)

		logRequest(log, statusCode, requestDuration, l.slowRequestThreshold, c.Path(), l.config)

		updateSpanWithResponseData(span, statusCode, c, requestDuration, err)

		return err
	}
}

func setSpanAttributes(span trace.Span, c *fiber.Ctx, requestID interface{}) {
	span.SetAttributes(
		attribute.String("http.method", c.Method()),
		attribute.String("http.url", c.OriginalURL()),
		attribute.String("http.host", c.Hostname()),
		attribute.String("http.user_agent", c.Get("User-Agent")),
		attribute.String("http.request_id", fmt.Sprintf("%v", requestID)),
	)
}

func extractResponseBody(c *fiber.Ctx) string {
	responseBody := c.Response().Body()
	if len(responseBody) == 0 {
		return ""
	}

	contentType := c.Response().Header.ContentType()

	if isTextContent(string(contentType)) {
		if len(responseBody) >= 2048 {
			return string(responseBody[:2048]) + "... (truncated)"
		}
		return string(responseBody)
	}

	return fmt.Sprintf("[binary content of type %s]", contentType)
}

func createLogEntry(
	c *fiber.Ctx,
	requestID interface{},
	statusCode int,
	duration time.Duration,
	responseBodyStr string,
	err error,
) *zap.SugaredLogger {
	log := logger.Log.With(
		"request_id", requestID,
		"method", c.Method(),
		"url", c.OriginalURL(),
		"path", c.Path(),
		"status", statusCode,
		"duration_ms", duration.Milliseconds(),
		"remote_addr", c.IP(),
		"x_forwarded_for", c.Get("X-Forwarded-For"),
		"host", c.Hostname(),
		"user_agent", c.Get("User-Agent"),
		"referer", c.Get("Referer"),
		"content_length", len(c.Body()),
		"accept", c.Get("Accept"),
		"content_type", c.Get("Content-Type"),
		"authorization", maskAuthorizationHeader(c.Get("Authorization")),
	)

	if user := c.Locals("user"); user != nil {
		log = log.With("user", fmt.Sprintf("%+v", user))
	}

	if responseBodyStr != "" {
		log = log.With("response_body", responseBodyStr)
	}

	if err != nil {
		log = log.With("error", err.Error())
	}

	return log
}

func (l *LoggingMiddleware) updateMetrics(c *fiber.Ctx, statusCode int, duration time.Duration) {
	routePath := getSafeRoutePath(c)
	statusStr := strconv.Itoa(statusCode)

	l.metrics.requestDuration.WithLabelValues(c.Method(), routePath, statusStr).Observe(duration.Seconds())
	l.metrics.requestSize.WithLabelValues(c.Method(), routePath).Observe(float64(len(c.Body())))
	l.metrics.responseSize.WithLabelValues(c.Method(), routePath, statusStr).Observe(float64(c.Response().Header.ContentLength()))
}

func logRequest(log *zap.SugaredLogger, statusCode int, duration time.Duration, slowThreshold time.Duration, path string, cfg *config.Config) {
	isError := statusCode >= 500
	isWarning := statusCode >= 400 && statusCode < 500
	isSlow := duration > slowThreshold

	switch {
	case isError:
		log.Error("http request failed with server error")
	case isWarning:
		log.Info("http request failed with client error")
	case isSlow:
		log.With("slow_request", true).Warn("slow http request")
	case cfg.PathForLevelInfo(path):
		log.Info("http request")
	default:
		log.Debug("http request")
	}
}

func updateSpanWithResponseData(span trace.Span, statusCode int, c *fiber.Ctx, duration time.Duration, err error) {
	span.SetAttributes(
		attribute.Int("http.status_code", statusCode),
		attribute.Int64("http.response_content_length", int64(c.Response().Header.ContentLength())),
		attribute.Float64("http.duration_ms", float64(duration.Milliseconds())),
	)

	if err != nil {
		span.RecordError(err)
	}
}

func getSafeRoutePath(c *fiber.Ctx) string {
	route := c.Route()
	if route == nil || route.Path == "" {
		return c.Path()
	}
	return route.Path
}

func maskAuthorizationHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	if len(authHeader) > 15 {
		return authHeader[:15] + "..."
	}
	return "****"
}

func isTextContent(contentType string) bool {
	for _, textType := range textContentTypes {
		if len(contentType) >= len(textType) && contentType[:len(textType)] == textType {
			return true
		}
	}
	return false
}
