package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/haadi-coder/bookmark-manager/internal/config"
	"github.com/haadi-coder/bookmark-manager/internal/server/handler"
	"github.com/haadi-coder/bookmark-manager/internal/storage"
)

const (
	RequestsLimit = 5
	LimitWindow   = time.Second
)

var (
	// AllowedOrigins: Using wildcard "*" is intentional to support local development
	// and browser extensions, whose origins are dynamic and cannot be predetermined
	// (e.g., Chrome/Firefox extension UUIDs change per installation).
	// While this weakens CORS security, the trade-off is acceptable in a local/development context.
	AllowedOrigins = []string{"*"}
	AllowedMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	ExposedHeaders = []string{"X-Total"}
)

type Server struct {
	*http.Server
}

func New(cfg config.Config, storage storage.Storage) *Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: AllowedOrigins,
		AllowedMethods: AllowedMethods,
		ExposedHeaders: ExposedHeaders,
	}))
	router.Use(httprate.Limit(
		RequestsLimit,
		LimitWindow,
		httprate.WithKeyFuncs(httprate.KeyByEndpoint),
		httprate.WithResponseHeaders(httprate.ResponseHeaders{
			Limit:      "X-RateLimit-Limit",
			Remaining:  "X-RateLimit-Remaining",
			Reset:      "X-RateLimit-Reset",
			RetryAfter: "Retry-After",
		}),
	))

	router.Get("/bookmarks", handler.GetBookmarks(context.Background(), storage))
	router.Post("/bookmarks", handler.CreateBookmark(context.Background(), storage))
	router.Patch("/bookmarks", handler.EditBookmark(context.Background(), storage))
	router.Delete("/bookmarks", handler.DeleteBookmark(context.Background(), storage))
	router.Get("/bookmarks/exists", handler.CheckBookmark(context.Background(), storage))
	router.Get("/bookmarks/export/html", handler.NetscapeBookmarks(context.Background(), storage))
	router.Get("/health", handler.CheckHealth())

	s := &http.Server{
		Addr:         cfg.Address(),
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
