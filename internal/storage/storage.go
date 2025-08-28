package storage

import (
	"context"
	"fmt"

	"github.com/haadi-coder/bookmark-manager/internal/model"
	"github.com/haadi-coder/bookmark-manager/internal/storage/postgresql"
)

const (
	Postgresql = "postgresql"
)

type Storage interface {
	GetBookmarks(ctx context.Context, limit, page int, query string) ([]*model.Bookmark, error)
	CreateBookmark(ctx context.Context, title, url string) (*model.Bookmark, error)
	EditBookmark(ctx context.Context, title, url string) (*model.Bookmark, error)
	DeleteBookmark(ctx context.Context, id int) error
	BookmarkExist(ctx context.Context, url string) (bool, error)
}

type Config struct {
	Type string
	Path string
}

func New(cfg *Config) (Storage, error) {
	switch cfg.Type {
	case Postgresql:
		return postgresql.New(cfg.Path)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
	}
}
