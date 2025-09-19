package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"github.com/haadi-coder/bookmark-manager/internal/api/request"
	"github.com/haadi-coder/bookmark-manager/internal/api/response"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"
	"github.com/haadi-coder/bookmark-manager/internal/model"
)

type BookmarkProvider interface {
	GetBookmarks(ctx context.Context, limit, offset int, search string) ([]*model.Bookmark, int, error)
}

func Bookmarks(ctx context.Context, provider BookmarkProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := request.ParseListOptions(r)
		if err != nil {
			slog.Error("failed to parse query params. Default params was applied", logger.Error(err))
		}

		result, totalCount, err := provider.GetBookmarks(ctx, opts.Perpage, opts.Offset(), opts.Search)
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
