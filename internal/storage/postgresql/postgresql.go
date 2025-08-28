package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/haadi-coder/bookmark-manager/internal/model"

	_ "github.com/lib/pq"
)

type PostgresqlStorage struct {
	db *sql.DB
}

func New(path string) (*PostgresqlStorage, error) {
	// TODO: использовать порт из конфига

	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	return &PostgresqlStorage{
		db: db,
	}, nil
}

func (s *PostgresqlStorage) GetBookmarks(ctx context.Context, limit, page int, query string) ([]*model.Bookmark, error) {
	return nil, nil
}

func (s *PostgresqlStorage) CreateBookmark(ctx context.Context, title, url string) (int64, error) {
	return 0, nil
}

func (s *PostgresqlStorage) EditBookmark(ctx context.Context, title, url string) (*model.Bookmark, error) {
	return nil, nil
}

func (s *PostgresqlStorage) DeleteBookmark(ctx context.Context, id string) error {
	return nil
}

func (s *PostgresqlStorage) BookmarkExist(ctx context.Context, url string) (bool, *model.Bookmark, error) {
	return false, nil, nil
}
