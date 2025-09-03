package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/haadi-coder/bookmark-manager/internal/config"
	"github.com/haadi-coder/bookmark-manager/internal/server/handler"
	"github.com/haadi-coder/bookmark-manager/internal/storage"
)

var (
	AllowedOrigins = []string{"http://localhost:5173"}
	AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
)

type Server struct {
	*http.Server
}

func New(cfg config.Config, storage storage.Storage) *Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: AllowedOrigins,
		AllowedMethods: AllowedMethods,
	}))

	router.Get("/bookmarks", handler.GetBookmarks(context.Background(), storage))
	router.Post("/bookmarks", handler.CreateBookmark(context.Background(), storage))
	router.Patch("/bookmarks", handler.EditBookmark(context.Background(), storage))
	router.Delete("/bookmarks", handler.DeleteBookmark(context.Background(), storage))
	router.Get("/bookmarks/exists", handler.CheckBookmark(context.Background(), storage))
	router.Get("/bookmarks/export/html", handler.NetscapeBookmarks(context.Background(), storage))

	s := &http.Server{
		Addr:         cfg.HttpAddress,
		Handler:      router,
		ReadTimeout:  cfg.HttpTimeout,
		WriteTimeout: cfg.HttpTimeout,
		IdleTimeout:  cfg.HttpIdleTimeout,
	}

	return &Server{
		Server: s,
	}
}

func (s *Server) Start() error {
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
