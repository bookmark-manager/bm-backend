package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/httprate"
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
	server *http.Server
}

func New(ctx context.Context, adress string, timeout, idleTimeout time.Duration, storage storage.Storage) *Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(httplog.RequestLogger(slog.Default(), &httplog.Options{
		Schema: &httplog.Schema{
			ErrorType:     "err_type",
			ErrorMessage:  "err_Msg",
			RequestBytes:  "req",
			ResponseBytes: "resp",
		},
		RecoverPanics: true,
	}))
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

	router.Get("/bookmarks", handler.GetBookmarks(ctx, storage))
	router.Post("/bookmarks", handler.CreateBookmark(ctx, storage))
	router.Patch("/bookmarks", handler.EditBookmark(ctx, storage))
	router.Delete("/bookmarks", handler.DeleteBookmark(ctx, storage))
	router.Get("/bookmarks/exists", handler.CheckBookmark(ctx, storage))
	router.Get("/bookmarks/export/html", handler.NetscapeBookmarks(ctx, storage))
	router.Get("/health", handler.CheckHealth())

	s := &http.Server{
		Addr:         adress,
		Handler:      router,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  idleTimeout,
	}

	return &Server{
		server: s,
	}
}

func (s *Server) Start() error {
	slog.Info("HTTP server starting", slog.String("address", s.server.Addr))

	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("HTTP server closed")
			return nil
		}
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("shutting down HTTP server")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	slog.Info("HTTP server shutdown completed")
	return nil
}
