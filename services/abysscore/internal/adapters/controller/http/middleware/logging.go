package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/config"
	"github.com/intezya/pkglib/logger"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
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
}

var textContentTypes = []string{
	"text/",
	"application/json",
	"application/xml",
	"application/javascript",
	"application/xhtml+xml",
}

const (
	defaultSlowRequestThreshold = 500 * time.Millisecond
	defaultSamplingRate         = 100
)

func NewLoggingMiddleware(config *config.Config) *LoggingMiddleware {
	metrics := initPrometheusMetrics()

	slowRequestThreshold := defaultSlowRequestThreshold
	if config.SlowRequestThresholdMs > 0 {
		slowRequestThreshold = time.Duration(config.SlowRequestThresholdMs) * time.Millisecond
	}

	return &LoggingMiddleware{
		config:               config,
		metrics:              metrics,
		resourceSampler:      newResourceSampler(defaultSamplingRate),
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
	const (
		bucketStart  = 100
		bucketFactor = 10
		bucketCount  = 8
	)

	metrics := &prometheusMetrics{
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct // useless (for this application) fields
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		requestSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct // useless (for this application) fields
				Name:    "http_request_size_bytes",
				Help:    "Size of HTTP requests in bytes",
				Buckets: prometheus.ExponentialBuckets(bucketStart, bucketFactor, bucketCount),
			},
			[]string{"method", "path"},
		),
		responseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct // useless (for this application) fields
				Name:    "http_response_size_bytes",
				Help:    "Size of HTTP responses in bytes",
				Buckets: prometheus.ExponentialBuckets(bucketStart, bucketFactor, bucketCount),
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

		logRequest(log, err, statusCode, requestDuration, l.slowRequestThreshold)

		updateSpanWithResponseData(span, statusCode, c, requestDuration, err)

		return err
	}
}

func setSpanAttributes(span trace.Span, c *fiber.Ctx, requestID interface{}) {
	span.SetAttributes(
		attribute.String("http.method", c.Method()),
		attribute.String("http.url", c.OriginalURL()),
		attribute.String("http.host", c.Hostname()),
		attribute.String("http.user_agent", c.Get("UserDTO-Agent")),
		attribute.String("http.request_id", fmt.Sprintf("%v", requestID)),
	)
}

func extractResponseBody(c *fiber.Ctx) string {
	const Large = 2048

	responseBody := c.Response().Body()
	if len(responseBody) == 0 {
		return ""
	}

	contentType := c.Response().Header.ContentType()

	if isTextContent(string(contentType)) {
		if len(responseBody) >= Large {
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
		"user_agent", c.Get("UserDTO-Agent"),
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

	l.metrics.requestDuration.WithLabelValues(c.Method(), routePath, statusStr).
		Observe(duration.Seconds())
	l.metrics.requestSize.WithLabelValues(c.Method(), routePath).Observe(float64(len(c.Body())))
	l.metrics.responseSize.WithLabelValues(
		c.Method(),
		routePath,
		statusStr,
	).Observe(float64(c.Response().Header.ContentLength()))
}

func logRequest(
	log *zap.SugaredLogger,
	err error,
	statusCode int,
	duration time.Duration,
	slowThreshold time.Duration,
) {
	isError := statusCode >= fiber.StatusInternalServerError
	isWarning := statusCode >= fiber.StatusBadRequest &&
		statusCode < fiber.StatusInternalServerError
	isSlow := duration > slowThreshold

	switch {
	case isError:
		log.Error("http request failed with server error")
	case isWarning:
		log.Info("http request failed with client error")
	case isSlow:
		log.With("slow_request", true).Warn("slow http request")
	default:
		log.Debug("http request")
	}
}

func updateSpanWithResponseData(
	span trace.Span,
	statusCode int,
	c *fiber.Ctx,
	duration time.Duration,
	err error,
) {
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
	const Long = 15

	if authHeader == "" {
		return ""
	}

	if len(authHeader) > Long {
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
