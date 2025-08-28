package main

import (
	"log/slog"
	"os"

	"github.com/haadi-coder/bookmark-manager/internal/config"
)

func main() {
	_ = config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("starting bookmark-manager")

	// ToDO: Storage

	// TODO: Router

	// TODO: Server run
}
