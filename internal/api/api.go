package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/httprate"
	"github.com/haadi-coder/bookmark-manager/internal/api/handler"
	"github.com/haadi-coder/bookmark-manager/internal/storage"
)

// Переименовать пакет в api
const (
	reqLimit  = 5
	reqWindow = time.Second
)

// Using wildcard "*" is intentional to support local development
// and browser extensions, whose origins are dynamic and cannot be predetermined
// (e.g., Chrome/Firefox extension UUIDs change per installation).
// While this weakens CORS security, the trade-off is acceptable in a local/development context.

type Server struct {
	server *http.Server
}

func NewServer(ctx context.Context, adress string, timeout, idleTimeout time.Duration, storage storage.Storage) *Server {
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
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		ExposedHeaders: []string{"X-Total"},
	}))
	router.Use(httprate.Limit(
		reqLimit,
		reqWindow,
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
	router.Get("/health", handler.CheckHealth(storage.Ping))

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

func (s *Server) Run(ctx context.Context) error {
	slog.Info("HTTP server starting", slog.String("address", s.server.Addr))

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- s.server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("failed to start http server: %w", err)
		}

	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown server gracefully: %w", err)
		}
	}

	return nil
}
