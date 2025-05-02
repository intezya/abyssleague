package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/intezya/abyssleague/services/websocket-messaging/cmd/app"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/config"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/controller/grpcapi"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/controller/websocket"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/infrastructure/hub"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/infrastructure/service"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/pkg/auth"
	"github.com/intezya/pkglib/logger"
)

const gracefulShutdownTimeout = 10 * time.Second

func main() {
	if _, err := os.Stat("static/websocket_debugger.html"); err != nil {
		fmt.Println("File not found:", err)
	} else {
		fmt.Println("File exists!")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	appConfig := config.Setup()
	jwtService := auth.NewJWTHelper(appConfig.JwtConfiguration())

	httpApp := app.NewHttpApp(appConfig)

	gRPCApps, hubs := setupHubs(ctx, httpApp, jwtService, appConfig)

	httpErrCh := startHTTPServer(ctx, httpApp)

	shutdownReason := waitForShutdownSignal(sigCh, httpErrCh)
	logger.Log.Infof("Shutting down application: %s", shutdownReason)

	gracefulShutdown(cancel, httpApp, gRPCApps, hubs)
}

func setupHubs(
	ctx context.Context,
	httpApp *app.HttpApp,
	jwtService *auth.JWTHelper,
	appConfig *config.Config,
) ([]*app.GRPCApp, []*hub.Hub) {
	gRPCApps := make([]*app.GRPCApp, 0, len(appConfig.Hubs))
	hubs := make([]*hub.Hub, 0, len(appConfig.Hubs))

	var wg sync.WaitGroup

	for idx, hubName := range appConfig.Hubs {
		logger.Log.Info("Starting hub: ", hubName)

		newHub := hub.NewHub(hubName)
		hubs = append(hubs, newHub)

		go newHub.Run()

		websocket.SetupRoute(httpApp.Mux, newHub, hubName, jwtService)

		websocketService := service.NewWebsocketService(newHub)
		gRPCApp := app.NewGRPCApp(appConfig.GRPCPorts[idx])
		grpcapi.Setup(gRPCApp.GRPCServer, websocketService)

		gRPCApps = append(gRPCApps, gRPCApp)

		wg.Add(1)

		go startGRPCServer(ctx, gRPCApp, &wg)
	}

	return gRPCApps, hubs
}

func startGRPCServer(ctx context.Context, app *app.GRPCApp, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := app.Start(ctx); err != nil {
		logger.Log.Errorf("gRPC server error: %v", err)
	}
}

func startHTTPServer(ctx context.Context, httpApp *app.HttpApp) chan error {
	httpErrCh := make(chan error, 1)
	go func() {
		if err := httpApp.Start(ctx); err != nil {
			httpErrCh <- err
		}

		close(httpErrCh)
	}()

	return httpErrCh
}

func waitForShutdownSignal(sigCh chan os.Signal, httpErrCh chan error) string {
	select {
	case <-sigCh:
		return "Received shutdown signal"
	case err := <-httpErrCh:
		if err != nil {
			return fmt.Sprintf("HTTP server error: %v", err)
		}

		return "HTTP server stopped unexpectedly"
	}
}

func gracefulShutdown(
	cancel context.CancelFunc,
	httpApp *app.HttpApp,
	gRPCApps []*app.GRPCApp,
	hubs []*hub.Hub,
) {
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		gracefulShutdownTimeout,
	)

	defer shutdownCancel()
	cancel()

	shutdownHTTPServer(shutdownCtx, httpApp)
	shutdownGRPCServers(shutdownCtx, gRPCApps)
	shutdownHubs(hubs)

	logger.Log.Info("Application shut down completed")
}

func shutdownHTTPServer(ctx context.Context, httpApp *app.HttpApp) {
	if err := httpApp.Shutdown(ctx); err != nil {
		logger.Log.Errorf("Error during HTTP server shutdown: %v", err)
	} else {
		logger.Log.Info("HTTP server shut down successfully")
	}
}

func shutdownGRPCServers(ctx context.Context, gRPCApps []*app.GRPCApp) {
	var wg sync.WaitGroup

	for _, grpcApp := range gRPCApps {
		wg.Add(1)

		go func(a *app.GRPCApp) {
			defer wg.Done()

			if err := a.Shutdown(ctx); err != nil {
				logger.Log.Errorf("Error during gRPC server shutdown: %v", err)
			}
		}(grpcApp)
	}

	wgDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	select {
	case <-wgDone:
		logger.Log.Info("All gRPC servers shut down gracefully")
	case <-ctx.Done():
		logger.Log.Warn("Shutdown timeout reached for gRPC servers, forcing exit")
	}
}

func shutdownHubs(hubs []*hub.Hub) {
	for _, h := range hubs {
		logger.Log.Infof("Hub %s stopped", h.GetName())
	}
}
