package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/haadi-coder/bookmark-manager/internal/api/response"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

func CheckHealth(pinger Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbStatus := "up"
		if err := pinger.Ping(r.Context()); err != nil {
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
