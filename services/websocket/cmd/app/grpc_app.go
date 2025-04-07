package app

import (
	"abysslib/logger"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type GRPCApp struct {
	GRPCServer *grpc.Server
	host       string
	port       int
	listener   net.Listener
}

func InterceptorLogger(l *zap.SugaredLogger) logging.Logger {
	return logging.LoggerFunc(
		func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			interceptorLogger := l.WithOptions(zap.AddCallerSkip(1)).With(fields...)
			switch lvl {
			case logging.LevelDebug:
				interceptorLogger.Debug(msg)
			case logging.LevelInfo:
				interceptorLogger.Info(msg)
			case logging.LevelWarn:
				interceptorLogger.Warn(msg)
			case logging.LevelError:
				interceptorLogger.Error(msg)
			default:
				panic(fmt.Sprintf("unknown level %v", lvl))
			}
		},
	)
}

func NewGRPCApp(host string, port int) *GRPCApp {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(
			func(p interface{}) (err error) {
				logger.Log.Errorf("Recovered from panic: %v", p)
				return status.Errorf(codes.Internal, "internal error")
			},
		),
	}

	kaParams := keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Minute,
		MaxConnectionAge:      30 * time.Minute,
		MaxConnectionAgeGrace: 5 * time.Minute,
		Time:                  5 * time.Minute,
		Timeout:               20 * time.Second,
	}

	kaPolicy := keepalive.EnforcementPolicy{
		MinTime:             1 * time.Minute,
		PermitWithoutStream: true,
	}

	GRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(InterceptorLogger(logger.Log.SugaredLogger), loggingOpts...),
		),
		grpc.KeepaliveParams(kaParams),
		grpc.KeepaliveEnforcementPolicy(kaPolicy),
	)

	return &GRPCApp{
		GRPCServer: GRPCServer,
		host:       host,
		port:       port,
	}
}

func (a *GRPCApp) Start(ctx context.Context) {
	logger.Log.Info("Starting gRPC server...")

	var err error
	a.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", a.host, a.port))
	if err != nil {
		logger.Log.Fatalf("Failed to listen on gRPC port: %v", err)
		return
	}

	logger.Log.Info("gRPC server listening on ", a.host, ":", a.port)

	errCh := make(chan error, 1)
	go func() {
		if err := a.GRPCServer.Serve(a.listener); err != nil {
			errCh <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		a.Stop()
	case err := <-errCh:
		logger.Log.Fatal("gRPC server failed: ", err)
	}
}

func (a *GRPCApp) Stop() {
	logger.Log.Info("Gracefully shutting down gRPC server...")

	stopped := make(chan struct{})
	go func() {
		a.GRPCServer.GracefulStop()
		close(stopped)
	}()

	timeout := time.After(10 * time.Second)
	select {
	case <-stopped:
		logger.Log.Info("gRPC server shutdown completed")
	case <-timeout:
		logger.Log.Warn("gRPC server shutdown timed out, forcing stop")
		a.GRPCServer.Stop()
	}
}
