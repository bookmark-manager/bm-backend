package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/haadi-coder/bookmark-manager/internal/api/response"
)

func CheckHealth(dbPing func(ctx context.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		dbStatus := "up"
		if err := dbPing(r.Context()); err != nil {
			dbStatus = "down"
		}

		if dbStatus == "down" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		render.JSON(w, r, response.HealthResponse{
			Status: dbStatus,
			Checks: response.HealthChecks{
				Postgres: dbStatus,
			},
		})
	}
}
