package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/haadi-coder/bookmark-manager/internal/api/response"
)

type BookmarkChecker interface {
	BookmarkExist(ctx context.Context, url string) (int, bool, error)
}

func CheckBookmark(ctx context.Context, checker BookmarkChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		url := r.URL.Query().Get("url")
		id, ok, err := checker.BookmarkExist(ctx, url)
		if err != nil {
			slog.Error("failed to check for bookmark", slog.String("url", url))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to check for bookmark"))
			return
		}

		if ok {
			slog.Info("bookmark with this url found", slog.String("url", url))
		} else {
			slog.Info("bookmark with this url not found", slog.String("url", url))
		}

		render.JSON(w, r, response.Response{
			Data: struct {
				ID    int  `json:"id"`
				Found bool `json:"found"`
			}{
				ID:    id,
				Found: ok,
			},
		})
	}
}
