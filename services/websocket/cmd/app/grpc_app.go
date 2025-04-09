package app

import (
	"abysslib/logger"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type GRPCApp struct {
	GRPCServer *grpc.Server
	port       int
	listener   net.Listener
}

func NewGRPCApp(port int) *GRPCApp {
	grpcServer := grpc.NewServer()

	return &GRPCApp{
		GRPCServer: grpcServer,
		port:       port,
	}
}

func (a *GRPCApp) Start(ctx context.Context) error {
	logger.Log.Infof("Starting gRPC server on port %d", a.port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", a.port, err)
	}
	a.listener = lis

	errCh := make(chan error, 1)
	go func() {
		if err := a.GRPCServer.Serve(lis); err != nil {
			errCh <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Log.Infof("Shutting down gRPC server on port %d...", a.port)
		return nil
	case err := <-errCh:
		return err
	}
}

func (a *GRPCApp) Shutdown(ctx context.Context) error {
	shutdownComplete := make(chan struct{})

	go func() {
		a.GRPCServer.GracefulStop()
		close(shutdownComplete)
	}()

	select {
	case <-shutdownComplete:
		logger.Log.Infof("gRPC server on port %d shutdown completed", a.port)
		return nil
	case <-ctx.Done():
		logger.Log.Warnf("gRPC server on port %d shutdown timed out, forcing stop", a.port)
		a.GRPCServer.Stop()
		return fmt.Errorf("shutdown timed out")
	}
}
