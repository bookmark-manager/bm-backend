package storage

import (
	"context"
	"errors"

	"github.com/haadi-coder/bookmark-manager/internal/model"
)

var (
	ErrNotFound = errors.New("bookmark not found")
	ErrExists   = errors.New("bookmark for this url already exists")
)

type Storage interface {
	GetBookmarks(ctx context.Context, limit, offset int, search string) ([]*model.Bookmark, int, error)
	CreateBookmark(ctx context.Context, title, url string) (*model.Bookmark, error)
	EditBookmark(ctx context.Context, id int, title, url string) (*model.Bookmark, error)
	DeleteBookmark(ctx context.Context, id int) error
	BookmarkExist(ctx context.Context, url string) (int, bool, error)
	Close() error
}
