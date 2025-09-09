package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/haadi-coder/bookmark-manager/internal/config"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"
	"github.com/haadi-coder/bookmark-manager/internal/server"
	"github.com/haadi-coder/bookmark-manager/internal/storage/postgresql"
)

func main() {
	slog.Info("starting bookmark-manager")

	if err := run(); err != nil {
		slog.Error("failed to start bookmark-manager", logger.Error(err))
		os.Exit(1)
	}

	slog.Info("server stopped")
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}

	setupLogger()

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

func setupLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
}
