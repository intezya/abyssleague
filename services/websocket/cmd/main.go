package main

import (
	"abysslib/dotenv"
	"context"
	"os"
	"os/signal"
	"syscall"
	"websocket/internal/infrastructure/hub"
	"websocket/internal/infrastructure/service"

	"abysslib/jwt"
	"websocket/cmd/app"
	"websocket/internal/adapters/config"
	"websocket/internal/adapters/controller/grpcapi"
	"websocket/internal/adapters/controller/ws"
)

type a struct{}

func (a) GetID() (id int)                { return 1 }
func (a) GetUsername() (username string) { return "0" }
func (a) GetHardwareID() (hwid string)   { return "0" }

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	dotenv.LoadEnv()
	appConfig := config.Configure()
	jwtService := jwt.New(appConfig)

	mainHub := hub.NewHub()
	go mainHub.Run()

	draftHub := hub.NewHub()
	go draftHub.Run()

	httpApp := app.NewHttpApp(appConfig)
	ws.SetupRoute(httpApp.Mux, mainHub, "main", jwtService)
	ws.SetupRoute(httpApp.Mux, draftHub, "draft", jwtService)
	httpAppDone := make(chan struct{})
	go func() {
		defer close(httpAppDone)
		httpApp.Start(ctx)
	}()

	mainWebsocketService := service.NewWebsocketService(mainHub)
	mainGRPCApp := app.NewGRPCApp(appConfig.MainGRPCPort)
	grpcapi.Setup(mainGRPCApp.GRPCServer, mainWebsocketService)
	mainGRPCAppDone := make(chan struct{})
	go func() {
		defer close(mainGRPCAppDone)
		mainGRPCApp.Start(ctx)
	}()

	draftWebsocketService := service.NewWebsocketService(draftHub)
	draftGRPCApp := app.NewGRPCApp(appConfig.DraftGRPCPort)
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
