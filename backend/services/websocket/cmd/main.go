package main

import (
	"abysslib/dotenv"
	"context"
	"os"
	"os/signal"
	"syscall"

	"abysslib/jwt"
	"websocket/cmd/app"
	"websocket/internal/adapters/config"
	"websocket/internal/adapters/controller/grpcapi"
	"websocket/internal/adapters/controller/ws"
	"websocket/internal/domain/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	dotenv.LoadEnv()
	appConfig := config.Configure()
	jwtService := jwt.New(appConfig)

	mainHub := service.NewHub()
	go mainHub.Run()
	ws.SetupRoute(mainHub, "main", jwtService)

	draftHub := service.NewHub()
	go draftHub.Run()
	ws.SetupRoute(draftHub, "draft", jwtService)

	httpApp := app.NewHttpApp(appConfig)
	httpAppDone := make(chan struct{})
	go func() {
		defer close(httpAppDone)
		httpApp.Start(ctx)
	}()

	mainWebsocketService := service.NewWebsocketService(mainHub)
	mainGRPCApp := app.NewGRPCApp(appConfig.MainGRPCHost, appConfig.MainGRPCPort)
	grpcapi.Setup(mainGRPCApp.GRPCServer, mainWebsocketService)
	mainGRPCAppDone := make(chan struct{})
	go func() {
		defer close(mainGRPCAppDone)
		mainGRPCApp.Start(ctx)
	}()

	draftWebsocketService := service.NewWebsocketService(draftHub)
	draftGRPCApp := app.NewGRPCApp(appConfig.DraftGRPCHost, appConfig.DraftGRPCPort)
	grpcapi.Setup(draftGRPCApp.GRPCServer, draftWebsocketService)
	draftGRPCAppDone := make(chan struct{})
	go func() {
		defer close(draftGRPCAppDone)
		draftGRPCApp.Start(ctx)
	}()

	select {
	case <-sigCh:
		cancel()
	}

	<-httpAppDone
	<-mainGRPCAppDone
	<-draftGRPCAppDone
}
