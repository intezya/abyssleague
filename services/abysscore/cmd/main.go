package main

import (
	"abysscore/internal/adapters/config"
	"abysscore/internal/adapters/controller/http/handlers"
	"abysscore/internal/adapters/controller/http/server"
	applicationservice "abysscore/internal/application/service"
	rediswrapper "abysscore/internal/infrastructure/cache/redis"
	"abysscore/internal/infrastructure/metrics/tracer"
	"abysscore/internal/infrastructure/persistence"
	"abysscore/pkg/auth"
	"abysscore/pkg/validator"
	_ "github.com/lib/pq"
	_ "net/http/pprof"
)

func main() {
	appConfig := config.LoadConfig()

	config.SetupLogger(appConfig.IsDebug, appConfig.EnvType, appConfig.LoggerConfig)

	tracerCleanup := tracer.Init(appConfig.TracerConfig)
	defer tracerCleanup()
	entClient := persistence.SetupEnt(appConfig.EntConfig)
	defer entClient.Close()
	redisClient := rediswrapper.NewClientWrapper(appConfig.RedisConfig)
	defer redisClient.Close()

	repositoryDependencies := persistence.NewDependencyProvider(entClient)

	serviceDependencies := applicationservice.NewDependencyProvider(
		repositoryDependencies,
		auth.NewHashHelper(),
		auth.NewJWTHelper(appConfig.JWTConfiguration),
	)

	handlerDependencies := handlers.NewDependencyProvider(
		serviceDependencies,
		validator.NewValidator(),
	)

	// TODO: provide redisClient.Client to handlers (check for pointer behavior)
	serverDependencies := server.NewDependencyProvider(
		redisClient,
		appConfig,
		handlerDependencies,
		serviceDependencies.AuthenticationService,
	)

	app := server.Setup(serverDependencies)

	server.Run(app)
}
