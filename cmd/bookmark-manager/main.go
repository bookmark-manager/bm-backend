package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/haadi-coder/bookmark-manager/internal/config"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"

	"github.com/haadi-coder/bookmark-manager/internal/storage"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("starting bookmark-manager")

	_, err := storage.New(&storage.Config{
		Type: storage.Postgresql,
		Path: fmt.Sprintf("postgres://%s:%s@localhost:5432/bookmarks?sslmode=disable", cfg.User, cfg.Password),
	})
	if err != nil {
		log.Error("failed to init storage", logger.Error(err))
		os.Exit(1)
	}

	// TODO: Router

	// TODO: Server run
}
