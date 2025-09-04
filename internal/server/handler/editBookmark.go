package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"
	"github.com/haadi-coder/bookmark-manager/internal/server/request"
	"github.com/haadi-coder/bookmark-manager/internal/server/response"
	"github.com/haadi-coder/bookmark-manager/internal/storage"
)

func EditBookmark(ctx context.Context, store storage.Storage) http.HandlerFunc {
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

		id := r.URL.Query().Get("id")
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			slog.Error("failed to get id from url", logger.Error(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		edited, err := store.EditBookmark(ctx, parsedId, reqData.Title, reqData.URL)
		if errors.Is(err, storage.ErrNotFound) {
			slog.Info(storage.ErrNotFound.Error(), slog.String("url", reqData.URL))

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error(storage.ErrNotFound.Error()))
			return
		}
		if err != nil {
			slog.Error("failed to edit bookmark", logger.Error(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to edit bookmark"))
			return
		}

		slog.Info("bookmark sucessfully edited", slog.Int("id", edited.ID))
		render.JSON(w, r, response.Response{
			Data: edited,
		})
	}
}
