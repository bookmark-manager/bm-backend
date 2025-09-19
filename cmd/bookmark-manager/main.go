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
	"github.com/haadi-coder/bookmark-manager/internal/storage"
	"github.com/lmittmann/tint"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	setupLogger(cfg)

	if err := run(cfg); err != nil {
		slog.Error("failed to start bookmark-manager", logger.Error(err))
		os.Exit(1)
	}
}

func run(cfg *config.Config) error {
	slog.Info("starting bookmark-manager")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	storage, err := storage.New(cfg.DB.DSN())
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	defer func() {
		if closeErr := storage.Close(); closeErr != nil {
			slog.Error("failed to close database connection", logger.Error(closeErr))
		}
	}()

	srv := api.NewServer(ctx, &api.ServerConfig{
		Address:     cfg.HTTP.Address(),
		Timeout:     cfg.HTTP.Timeout,
		IdleTimeout: cfg.HTTP.IdleTimeout,

		BookmarkProvider: storage,
		BookmarkChecker:  storage,
		BookmarkDeleter:  storage,
		BookmarkPinger:   storage,
		BookmarkEditor:   storage,
		BookmarkCreator:  storage,
	})

	if err := srv.Run(ctx); err != nil {
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
