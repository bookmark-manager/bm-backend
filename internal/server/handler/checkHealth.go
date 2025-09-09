package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/haadi-coder/bookmark-manager/internal/server/response"
)

func CheckHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		render.JSON(w, r, response.Response{
			Data: "service is healthy",
		})
	}
}
