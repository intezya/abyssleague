package main

import (
	"context"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/config"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/factory"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/wrapper"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/handlers"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/server"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/server/routes"
	applicationservice "github.com/intezya/abyssleague/services/abysscore/internal/application/service"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/mail"
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

	tracerCleanup := tracer.Init(appConfig.TracerConfig, logger.Log)
	entClient := persistence.SetupEnt(appConfig.EntConfig, logger.Log)
	redisClient := rediswrapper.NewClientWrapper(appConfig.RedisConfig, logger.Log)
	grpcFactory := factory.NewGrpcClientFactory()
	smtpClient := mail.NewSMTPSender(appConfig.SMTPConfig, logger.Log)

	logger.Log.Debug("grpcFactory has been initialized")

	defer func() {
		tracerCleanup()
		redisClient.Close()
		grpcFactory.CloseAll()

		_ = entClient.Close()
	}()

	gRPCDependencies := wrapper.NewDependencyProvider(
		context.Background(),
		appConfig.GRPCConfig,
		grpcFactory,
	)

	logger.Log.Debug("grpcDependencies has been initialized")

	repositoryDependencies := persistence.NewDependencyProvider(entClient, redisClient)

	serviceDependencies := applicationservice.NewDependencyProvider(
		repositoryDependencies,
		gRPCDependencies,
		auth.NewHashHelper(),
		auth.NewJWTHelper(appConfig.JWTConfiguration),
		smtpClient,
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
