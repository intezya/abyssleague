package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/config"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/controller/http/middleware"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/controller/http/routes"
	"github.com/intezya/pkglib/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	maxHeaderBytes    = 1 << 20 // 1 mb
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	readHeaderTimeout = 10 * time.Second
)

type HttpApp struct {
	Mux     *http.ServeMux
	config  *config.Config
	Server  *http.Server
	running bool
}

func NewHttpApp(config *config.Config) *HttpApp {
	mux := http.NewServeMux()

	mux.HandleFunc(
		routes.PingPath, func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("pong!"))
		},
	)

	mux.HandleFunc(routes.MetricsPath, promhttp.Handler().ServeHTTP)

	loggedHandler := middleware.RequestIDMiddleware(middleware.LoggingMiddleware(mux))

	server := &http.Server{
		Addr:                         fmt.Sprintf(":%d", config.HTTPPort),
		Handler:                      loggedHandler,
		ReadTimeout:                  readTimeout,
		WriteTimeout:                 writeTimeout,
		IdleTimeout:                  idleTimeout,
		ReadHeaderTimeout:            readHeaderTimeout,
		MaxHeaderBytes:               maxHeaderBytes, // 1 MB
		TLSConfig:                    nil,
		DisableGeneralOptionsHandler: false,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
		HTTP2:                        nil,
		Protocols:                    nil,
	}

	return &HttpApp{
		Mux:     mux,
		config:  config,
		Server:  server,
		running: false,
	}
}

func (a *HttpApp) Start(ctx context.Context) error {
	logger.Log.Infof("HTTP Server starting on port %d", a.config.HTTPPort)
	a.running = true

	errCh := make(chan error, 1)
	go func() {
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("HTTP server error: %w", err)
		}

		close(errCh)
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info("HTTP server context cancelled")

		return nil
	case err := <-errCh:
		a.running = false

		if err != nil {
			logger.Log.Errorf("HTTP server failed: %v", err)

			return err
		}

		return nil
	}
}

func (a *HttpApp) Shutdown(ctx context.Context) error {
	if !a.running {
		logger.Log.Info("HTTP server is not running, nothing to shutdown")

		return nil
	}

	logger.Log.Info("Shutting down HTTP server...")

	if err := a.Server.Shutdown(ctx); err != nil {
		logger.Log.Errorf("HTTP server shutdown error: %v", err)

		return err
	}

	a.running = false

	logger.Log.Info("HTTP server shutdown completed")

	return nil
}
