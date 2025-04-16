package server

import (
	"abysscore/internal/adapters/controller/http/middleware"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/intezya/pkglib/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func Setup(dependencies *DependencyProvider) *fiber.App {
	config := dependencies.config

	// Setup /metrics on separated http server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Log.Warn(http.ListenAndServe(fmt.Sprintf(":%d", config.MetricsPort), nil))
	}()

	server := fiber.New(
		fiber.Config{
			Prefork:                      false, // multicore support for performance
			StrictRouting:                true,
			CaseSensitive:                true,
			BodyLimit:                    10 * 1024 * 1024, // 10 MB
			Concurrency:                  100,              // max concurrent connections (requests)
			ReadTimeout:                  5 * time.Second,
			WriteTimeout:                 5 * time.Second,
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

	rateLimitMiddleware := middleware.NewRateLimitMiddleware(
		dependencies.redisClient,
		config,
	)
	recoverMiddleware := middleware.NewRecoverMiddleware(config.FiberRequestIDConfig)
	authenticationMiddleware := middleware.NewAuthenticationMiddleware(
		config.UnprotectedAuthRequests,
		dependencies.authenticationService,
		dependencies.redisClient,
	)
	loggingMiddleware := middleware.NewLoggingMiddleware(config)

	if config.IsDebug {
		logger.Log.Info("Setting up pprof middleware with endpoint: ", dependencies.config.FiberPprofConfig.Prefix)

		server.Use(pprof.New())
	}

	server.Use(requestid.New(config.FiberRequestIDConfig))
	server.Use(healthcheck.New(config.FiberHealthCheckConfig))
	server.Use(loggingMiddleware.Handle())
	server.Use(recoverMiddleware.Handle())
	server.Use(rateLimitMiddleware.HandleForAuth())
	server.Use(rateLimitMiddleware.HandleDefault())
	server.Use(authenticationMiddleware.Handle())

	apiGroup := server.Group("/api")

	authGroup := apiGroup.Group(config.Paths.Authentication.Self)
	authGroup.Post(config.Paths.Authentication.Register, dependencies.authenticationHandler.Register)
	authGroup.Post(config.Paths.Authentication.Login, dependencies.authenticationHandler.Login)

	return server
}

func Run(server *fiber.App) {
	log.Fatal(server.Listen(":8080"))
}
