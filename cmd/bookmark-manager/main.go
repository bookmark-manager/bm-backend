package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/haadi-coder/bookmark-manager/internal/config"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"
	"github.com/haadi-coder/bookmark-manager/internal/server"
	"github.com/haadi-coder/bookmark-manager/internal/storage"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(log)

	slog.Info("starting bookmark-manager")

	storage, err := storage.New(&storage.Config{
		Type: storage.Postgresql,
		Path: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		),
	})
	if err != nil {
		slog.Error("failed to init storage", logger.Error(err))
		return fmt.Errorf("failed to init storage: %w", err)
	}

	server := server.New(*cfg, storage)

	slog.Info("starting HTTP server", slog.String("address", server.Addr))

	if err = server.Start(); err != nil {
		slog.Error("failed to start server", logger.Error(err))
		return err
	}

	log.Info("server stopped")

	return nil
}
