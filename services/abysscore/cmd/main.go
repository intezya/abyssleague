package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/config"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/clients"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/handlers"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/server"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/server/routes"
	applicationservice "github.com/intezya/abyssleague/services/abysscore/internal/application/service"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/mail"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/persistence"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/auth"
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
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

	config.SetupLogger(appConfig.IsDebug, string(appConfig.EnvType), appConfig.LoggerConfig)

	tracerCleanup := tracer.Init(appConfig.TracerConfig, logger.Log)
	entClient := persistence.SetupEnt(appConfig.EntConfig, logger.Log)
	redisClient := rediswrapper.NewClientWrapper(appConfig.RedisConfig, logger.Log)
	smtpClient := mail.NewSMTPSender(appConfig.SMTPConfig, logger.Log)
	gRPCDependencies := clients.NewDependencyProvider(appConfig.GRPCConfig)

	defer func() {
		tracerCleanup()
		redisClient.Close()

		_ = entClient.Close()
		_ = gRPCDependencies.CloseAll()
	}()

	logger.Log.Debug("grpcDependencies has been initialized")

	repositoryDependencies := persistence.NewDependencyProvider(entClient, redisClient)

	serviceDependencies := applicationservice.NewDependencyProvider(
		repositoryDependencies,
		gRPCDependencies,
		auth.NewHashHelper(appConfig.HardwareIDEncryptionKey),
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

	errorz.SetValidator(validator.New())

	app := server.Setup(serverDependencies)

	server.Run(app, appConfig)
}
