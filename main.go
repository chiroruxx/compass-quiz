// Package main runs the Compass Quiz HTTP server.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:generate go tool oapi-codegen -config oapi-codegen.yaml openapi.yaml

const (
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 5 * time.Second
)

func main() {
	err := run()
	if err != nil {
		slog.Error("HTTP server stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = newRouter()
	server.ReadHeaderTimeout = readHeaderTimeout

	serverError := make(chan error, 1)
	go func() {
		serverError <- server.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
	}

	shutdownContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(shutdownContext)
	if err != nil {
		return fmt.Errorf("shut down HTTP server: %w", err)
	}

	slog.Info("HTTP server shut down")

	return nil
}
