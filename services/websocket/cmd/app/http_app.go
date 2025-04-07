package app

import (
	"abysslib/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"websocket/internal/adapters/config"
)

type HttpApp struct {
	config *config.Config
	server *http.Server
}

func NewHttpApp(config *config.Config) *HttpApp {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.HTTPPort),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	includePingHandler()

	return &HttpApp{
		config: config,
		server: server,
	}
}

func includePingHandler() {
	http.HandleFunc(
		"/ping", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("pong!"))
		},
	)

}

func (a *HttpApp) Start(ctx context.Context) {
	logger.Log.Infof("HTTP Server starting on port %d", a.config.HTTPPort)

	errCh := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			logger.Log.Errorf("HTTP server shutdown error: %v", err)
		} else {
			logger.Log.Info("HTTP server shutdown completed")
		}
	case err := <-errCh:
		logger.Log.Fatalf("HTTP server failed: %v", err)
	}
}
