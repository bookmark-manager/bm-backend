package handler

import (
	"context"
	"log/slog"
	"math"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/haadi-coder/bookmark-manager/internal/api/response"
	"github.com/haadi-coder/bookmark-manager/internal/lib/logger"
	"github.com/haadi-coder/bookmark-manager/internal/storage"
	"github.com/virtualtam/netscape-go"
	"github.com/virtualtam/netscape-go/types"
)

func NetscapeBookmarks(ctx context.Context, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		result, totalCount, err := storage.GetBookmarks(ctx, math.MaxInt32, 0, "")
		if err != nil {
			slog.Error("failed to get bookmarks from db", logger.Error(err))

			render.JSON(w, r, response.Error("failed to get bookmarks"))
			return
		}

		doc := types.Document{
			Title: "Bookmarks",
			Root: types.Folder{
				Name:      "Bookmarks",
				Bookmarks: make([]types.Bookmark, 0, totalCount),
			},
		}

		for _, bm := range result {
			doc.Root.Bookmarks = append(doc.Root.Bookmarks, types.Bookmark{
				Title:     bm.Title,
				CreatedAt: &bm.CreatedAt,
				UpdatedAt: &bm.UpdatedAt,
				Href:      bm.URL,
			})
		}

		m, err := netscape.Marshal(&doc)
		if err != nil {
			slog.Error("failed to marshal bookmarks to netscape format", logger.Error(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to marshal bookmarks to netscape format"))
			return
		}

		slog.Info("bookmarks successfully exported to netscape format",
			slog.Int("bookmarks_count", len(result)),
			slog.Int("output_size_bytes", len(m)))

		w.Header().Set("Content-Disposition", `attachment; filename="bookmarks.html"`)
		render.SetContentType(render.ContentTypePlainText)

		render.Data(w, r, m)
	}
}
