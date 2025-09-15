package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/haadi-coder/bookmark-manager/internal/api/request"
	"github.com/haadi-coder/bookmark-manager/internal/api/response"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"
	"github.com/haadi-coder/bookmark-manager/internal/storage"
)

func GetBookmarks(ctx context.Context, store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		opts, err := request.ParseListOptions(r)
		if err != nil {
			slog.Error("failed to parse query params. Default params was applied", logger.Error(err))
		}

		result, totalCount, err := store.GetBookmarks(ctx, opts.Perpage, opts.Offset(), opts.Search)
		if err != nil {
			slog.Error("failed to get bookmarks from db", logger.Error(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get bookmarks"))
			return
		}

		slog.Info("got bookmarks", slog.Any("bookmarks count", len(result)))

		w.Header().Set("X-Total", strconv.Itoa(totalCount))
		render.JSON(w, r, response.Response{
			Data: result,
		})
	}
}
