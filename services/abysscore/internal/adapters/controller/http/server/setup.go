package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	_ "github.com/intezya/abyssleague/services/abysscore/docs" // Path to generated docs
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/config"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/middleware"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/server/routes"
	"github.com/intezya/pkglib/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	bodyLimit                = 10 * 1024 * 1024 // 10 mb
	maxConcurrentConnections = 100
	writeTimeout             = 5 * time.Second
	readTimeout              = 5 * time.Second
)

// setupMetricsServer creates a separate HTTP server for Prometheus metrics.
func setupMetricsServer(port int) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Log.Warn(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()
}

// createFiberApp creates and configures a new Fiber application.
func createFiberApp(config *config.Config) *fiber.App {
	return fiber.New(
		//nolint:exhaustruct // used default if not changed
		fiber.Config{
			Prefork:                      false,
			StrictRouting:                true,
			CaseSensitive:                true,
			BodyLimit:                    bodyLimit,
			Concurrency:                  maxConcurrentConnections,
			ReadTimeout:                  readTimeout,
			WriteTimeout:                 writeTimeout,
			DisableKeepalive:             false,
			DisableDefaultContentType:    false,
			DisablePreParseMultipartForm: true,
			ReduceMemoryUsage:            false,
			JSONEncoder:                  jsoniter.Marshal,
			JSONDecoder:                  jsoniter.Unmarshal,
			EnablePrintRoutes:            config.IsDebug,
			DisableStartupMessage:        !config.IsDebug,
		},
	)
}

// setupCoreMiddleware sets up the common middleware for all routes.
func setupCoreMiddleware(app *fiber.App, config *config.Config) {
	if config.IsDebug {
		// TODO: add debug logging middleware
		logger.Log.Debug("Setting up pprof middleware")
		app.Use(pprof.New())
		logger.Log.Debugw("Setting up swagger middleware")
		app.Get("/swagger/**", swagger.HandlerDefault)
	}

	app.Use(requestid.New(config.FiberRequestIDConfig))
	app.Use(healthcheck.New(config.FiberHealthCheckConfig))
}

// createMiddlewareLinker creates all application middleware and links them.
func createMiddlewareLinker(
	dependencies *routes.DependencyProvider,
	config *config.Config,
) *routes.MiddlewareLinker {
	loggingMiddleware := middleware.NewLoggingMiddleware(config)
	recoverMiddleware := middleware.NewRecoverMiddleware(config.FiberRequestIDConfig)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(
		dependencies.RedisClient,
		config,
	)
	authenticationMiddleware := middleware.NewAuthenticationMiddleware(
		dependencies.AuthenticationService,
		dependencies.RedisClient,
	)

	return routes.NewMiddlewareLinker(
		loggingMiddleware,
		recoverMiddleware,
		rateLimitMiddleware,
		authenticationMiddleware,
	)
}

func Setup(dependencies *routes.DependencyProvider) *fiber.App {
	config := dependencies.Config

	// Set up metrics server on separate port
	setupMetricsServer(config.MetricsPort)

	// Create and configure main Fiber app
	server := createFiberApp(config)

	// Set up core middleware
	setupCoreMiddleware(server, config)

	// Create custom middleware linker
	middlewareLinker := createMiddlewareLinker(dependencies, config)

	// Register routes using route groups
	for _, group := range dependencies.GetRouteGroups() {
		group.Register(server, middlewareLinker)
	}

	return server
}

// Run starts the server with graceful shutdown support.
func Run(server *fiber.App, config *config.Config) {
	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Log.Info("Gracefully shutting down...")

		_ = server.Shutdown()
	}()

	// Start server
	port := fmt.Sprintf(":%d", config.ServerPort)
	logger.Log.Infof("Starting server on port %s", port)

	if err := server.Listen(port); err != nil {
		logger.Log.Fatalf("Server error: %v", err)
	}
}
