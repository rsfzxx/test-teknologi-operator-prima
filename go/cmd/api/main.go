package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/config"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/database"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/handler"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/repository/postgres"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/router"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application stopped with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	slog.Info("connected to PostgreSQL")

	if err := database.Migrate(ctx, pool); err != nil {
		return err
	}

	bankRepo := postgres.NewBankRepository(pool)
	bankService := service.NewBankService(bankRepo)

	apiHandler := router.New(router.Dependencies{
		Health: handler.NewHealthHandler(pool),
		Bank:   handler.NewBankHandler(bankService),
	})

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           apiHandler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("server started", "addr", srv.Addr)
		serverErr <- srv.ListenAndServe()
	}()

	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-sigCtx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	slog.Info("server stopped gracefully")
	return nil
}