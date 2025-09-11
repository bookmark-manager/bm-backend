package main

import (
	"fmt"
	"log/slog"
	"os"
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
	}

	setupLogger(cfg)

	slog.Info("starting bookmark-manager")

	if err := run(cfg); err != nil {
		slog.Error("failed to start bookmark-manager", logger.Error(err))
		os.Exit(1)
	}

	slog.Info("server stopped")
}

func run(cfg *config.Config) error {
	storage, err := postgresql.New(cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	server := server.New(*cfg, storage)
	if err = server.Start(); err != nil {
		return err
	}

	return nil
}

func setupLogger(cfg *config.Config) {
	w := os.Stdout
	logger := slog.New(tint.NewHandler(
		w,
		&tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			NoColor:    cfg.NoColor,
		},
	))
	slog.SetDefault(logger)
}
