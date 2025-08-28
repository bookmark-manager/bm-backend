package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/haadi-coder/bookmark-manager/internal/model"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const UniqueViolation = "23505"

var (
	ErrNotFound = errors.New("bookmark not found")
	ErrExists   = errors.New("bookmark already exists")
)

type PostgresqlStorage struct {
	db *sql.DB
}

func New(path string) (*PostgresqlStorage, error) {
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	return &PostgresqlStorage{
		db: db,
	}, nil
}

func (s *PostgresqlStorage) GetBookmarks(ctx context.Context, limit, page int, query string) ([]*model.Bookmark, error) {
	if limit == 0 {
		limit = 20
	}

	if page == 0 {
		page = 1
	}

	var rows *sql.Rows
	var err error

	offset := (page - 1) * limit

	if query != "" {
		rows, err = s.db.QueryContext(ctx, "SELECT * FROM bookmarks WHERE url ILIKE $1 OR title ILIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", query, limit, offset)
	} else {
		rows, err = s.db.QueryContext(ctx, "SELECT * FROM bookmarks ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarks rows: %w", err)
	}

	var bookmarks []*model.Bookmark
	for rows.Next() {
		var bm model.Bookmark

		if err := rows.Scan(&bm.ID, &bm.Title, &bm.URL, &bm.CreatedAt, &bm.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan bookmark: %w", err)
		}

		bookmarks = append(bookmarks, &bm)
	}

	return bookmarks, nil
}

func (s *PostgresqlStorage) CreateBookmark(ctx context.Context, title, url string) (*model.Bookmark, error) {
	var bm model.Bookmark

	row := s.db.QueryRowContext(ctx, "INSERT INTO bookmarks (title, url) VALUES($1, $2) RETURNING id, url, title, created_at, updated_at", title, url)

	err := row.Scan(&bm.ID, &bm.Title, &bm.URL, &bm.CreatedAt, &bm.UpdatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == UniqueViolation {
			return nil, ErrExists
		}

		return nil, fmt.Errorf("failed to create bookmark: %w", err)
	}

	return &bm, nil
}

func (s *PostgresqlStorage) EditBookmark(ctx context.Context, title, url string) (*model.Bookmark, error) {
	var bm model.Bookmark

	row := s.db.QueryRowContext(ctx, "UPDATE bookmarks SET title=$1, url=$2 WHERE url=$2 RETURNING id, url, title, created_at, updated_at", title, url)

	err := row.Scan(&bm.ID, &bm.Title, &bm.URL, &bm.CreatedAt, &bm.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to edit bookmark: %w", err)
	}

	return &bm, nil
}

func (s *PostgresqlStorage) DeleteBookmark(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM bookmarks WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get check deletion: %w", err)
	}

	if rowAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresqlStorage) BookmarkExist(ctx context.Context, url string) (bool, error) {
	var found bool

	err := s.db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM bookmarks WHERE url=$1)", url).Scan(&found)
	if err != nil {
		return false, fmt.Errorf("failed to find bookmark: %w", err)
	}

	return found, nil
}
