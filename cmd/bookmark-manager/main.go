package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/haadi-coder/bookmark-manager/internal/api"
	"github.com/haadi-coder/bookmark-manager/internal/config"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"
	"github.com/haadi-coder/bookmark-manager/internal/storage/postgres"

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
	storage, err := postgres.New(cfg.DB.DSN())
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	defer func() {
		if closeErr := storage.Close(); closeErr != nil {
			slog.Error("failed to close database connection", logger.Error(closeErr))
		}
	}()

	server := api.NewServer(ctx, cfg.Http.Address(), cfg.Http.Timeout, cfg.Http.IdleTimeout, storage)

	if err := server.Run(ctx); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}

func setupLogger(cfg *config.Config) {
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}

	logger := slog.New(tint.NewHandler(
		os.Stdout,
		&tint.Options{
			Level:      level,
			TimeFormat: time.Kitchen,
			NoColor:    cfg.NoColor,
		},
	))
	slog.SetDefault(logger)
}
