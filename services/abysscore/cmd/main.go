package main

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/config"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/factory"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/wrapper"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/handlers"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/server"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/server/routes"
	applicationservice "github.com/intezya/abyssleague/services/abysscore/internal/application/service"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/persistence"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/auth"
	"github.com/intezya/pkglib/logger"
	_ "github.com/lib/pq"
)

// @title						AbyssCore API
// @version					1.0
// @description				API of AbyssCore server
// @host						localhost:8080
// @BasePath					/
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @schemes					http https.
func main() {
	appConfig := config.LoadConfig()

	config.SetupLogger(appConfig.IsDebug, appConfig.EnvType, appConfig.LoggerConfig)

	tracerCleanup := tracer.Init(appConfig.TracerConfig)
	entClient := persistence.SetupEnt(appConfig.EntConfig)
	redisClient := rediswrapper.NewClientWrapper(appConfig.RedisConfig)
	grpcFactory := factory.NewGrpcClientFactory()

	logger.Log.Debug("grpcFactory has been initialized")

	defer func() {
		tracerCleanup()
		redisClient.Close()
		grpcFactory.CloseAll()

		_ = entClient.Close()
	}()

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
