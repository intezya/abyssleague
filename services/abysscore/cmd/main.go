package main

import (
	"abysscore/internal/adapters/config"
	"abysscore/internal/adapters/controller/grpc/factory"
	"abysscore/internal/adapters/controller/grpc/wrapper"
	"abysscore/internal/adapters/controller/http/handlers"
	"abysscore/internal/adapters/controller/http/server"
	"abysscore/internal/adapters/controller/http/server/routes"
	applicationservice "abysscore/internal/application/service"
	rediswrapper "abysscore/internal/infrastructure/cache/redis"
	"abysscore/internal/infrastructure/metrics/tracer"
	"abysscore/internal/infrastructure/persistence"
	"abysscore/internal/pkg/auth"
	"github.com/intezya/pkglib/logger"
	_ "github.com/lib/pq"
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

	grpcFactory := factory.NewGrpcClientFactory()
	defer grpcFactory.CloseAll()
	logger.Log.Debug("grpcFactory has been initialized")

	gRPCDependencies := wrapper.NewDependencyProvider(appConfig.GRPCConfig, grpcFactory)
	logger.Log.Debug("grpcDependencies has been initialized")

	repositoryDependencies := persistence.NewDependencyProvider(entClient)

	serviceDependencies := applicationservice.NewDependencyProvider(
		repositoryDependencies,
		gRPCDependencies,
		auth.NewHashHelper(),
		auth.NewJWTHelper(appConfig.JWTConfiguration),
	)

	handlerDependencies := handlers.NewDependencyProvider(serviceDependencies)

	serverDependencies := routes.NewDependencyProvider(
		handlerDependencies,
		serviceDependencies.AuthenticationService,
		redisClient,
		appConfig,
	)

	app := server.Setup(serverDependencies)

	server.Run(app, appConfig)
}
