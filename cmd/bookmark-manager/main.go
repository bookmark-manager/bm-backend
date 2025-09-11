package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/haadi-coder/bookmark-manager/internal/config"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"
	"github.com/haadi-coder/bookmark-manager/internal/server"
	"github.com/haadi-coder/bookmark-manager/internal/storage/postgresql"
	"github.com/lmittmann/tint"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to initialize application: %w", err)
		os.Exit(1)
	}

	setupLogger(cfg)

	slog.Info("starting bookmark-manager")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, cfg); err != nil {
		slog.Error("failed to start bookmark-manager", logger.Error(err))
		os.Exit(1)
	}

	slog.Info("server stopped")
}

func run(ctx context.Context, cfg *config.Config) error {
	storage, err := postgresql.New(cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	defer func() {
		slog.Info("closing database connection")
		if closeErr := storage.Close(); closeErr != nil {
			slog.Error("failed to close database connection", logger.Error(closeErr))
		}
	}()

	server := server.New(ctx, cfg.Address(), cfg.Http.Timeout, cfg.Http.IdleTimeout, storage)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Start()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("failed to start http server: %w", err)
		}

	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown server gracefully: %w", err)
		}
	}

	return nil
}

func setupLogger(cfg *config.Config) {
	w := os.Stdout

	level := slog.LevelInfo
	if cfg.Logger.Debug {
		level = slog.LevelDebug
	}

	logger := slog.New(tint.NewHandler(
		w,
		&tint.Options{
			Level:      level,
			TimeFormat: time.Kitchen,
			NoColor:    cfg.Logger.NoColor,
		},
	))
	slog.SetDefault(logger)
}
