package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"
	"github.com/haadi-coder/bookmark-manager/internal/server/response"
	"github.com/haadi-coder/bookmark-manager/internal/storage"
	"github.com/haadi-coder/bookmark-manager/internal/storage/postgresql"
)

func DeleteBookmark(ctx context.Context, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := r.URL.Query().Get("id")
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			slog.Error("failed to convert limit to integer", logger.Error(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		err = storage.DeleteBookmark(ctx, parsedId)
		if errors.Is(err, postgresql.ErrNotFound) {
			slog.Info(postgresql.ErrNotFound.Error(), slog.String("id", id))

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error(postgresql.ErrNotFound.Error()))
			return
		}
		if err != nil {
			slog.Error("failed to delete bookmark", logger.Error(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete bookmark"))
			return
		}

		slog.Info("bookmark sucessfully deleted", slog.String("id", id))
		render.JSON(w, r, response.Response{
			Data: "bookmark sucessfully deleted",
		})
	}
}
