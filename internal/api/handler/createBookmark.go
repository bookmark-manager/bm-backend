package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/haadi-coder/bookmark-manager/internal/api/request"
	"github.com/haadi-coder/bookmark-manager/internal/api/response"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"
	"github.com/haadi-coder/bookmark-manager/internal/storage"
)

func CreateBookmark(ctx context.Context, store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var reqData request.Request
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			slog.Error("failed to decode request body", logger.Error(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request body"))
			return
		}

		if err := validator.New().Struct(reqData); err != nil {
			slog.Error("invalid request", logger.Error(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		new, err := store.CreateBookmark(ctx, reqData.Title, reqData.URL)
		if errors.Is(err, storage.ErrExists) {
			slog.Info(storage.ErrExists.Error(), slog.String("url", reqData.URL))

			render.Status(r, http.StatusConflict)
			render.JSON(w, r, response.Error(storage.ErrExists.Error()))
			return
		}
		if err != nil {
			slog.Error("failed to create bookmark", logger.Error(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create bookmark"))
			return
		}

		slog.Info("bookmark sucessfully created", slog.Int("id", new.ID))
		render.JSON(w, r, response.Response{
			Data: new,
		})
	}
}
