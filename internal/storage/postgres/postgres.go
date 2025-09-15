package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/haadi-coder/bookmark-manager/internal/model"
	"github.com/haadi-coder/bookmark-manager/internal/storage"

	"github.com/lib/pq"
)

const uniqueViolation = "23505"

type PostgresStorage struct {
	db *sql.DB
}

func New(path url.URL) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", path.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) GetBookmarks(ctx context.Context, limit, offset int, search string) ([]*model.Bookmark, int, error) {
	var totalCount int
	var countErr error

	var rows *sql.Rows
	var queryErr error

	if search != "" {
		countErr = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM bookmarks WHERE url ILIKE '%' || $1 || '%' OR title ILIKE '%' || $1 || '%'", search).Scan(&totalCount)
		rows, queryErr = s.db.QueryContext(ctx, "SELECT * FROM bookmarks WHERE url ILIKE '%' || $1 || '%'  OR title ILIKE '%' || $1 || '%'  ORDER BY created_at DESC LIMIT $2 OFFSET $3", search, limit, offset)
	} else {
		countErr = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM bookmarks").Scan(&totalCount)
		rows, queryErr = s.db.QueryContext(ctx, "SELECT * FROM bookmarks ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	}

	if queryErr != nil {
		return nil, 0, fmt.Errorf("failed to get bookmarks rows: %w", queryErr)
	}

	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", countErr)
	}

	bookmarks := []*model.Bookmark{}
	for rows.Next() {
		var bm model.Bookmark

		if err := rows.Scan(&bm.ID, &bm.URL, &bm.Title, &bm.CreatedAt, &bm.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan bookmark: %w", err)
		}

		bookmarks = append(bookmarks, &bm)
	}

	return bookmarks, totalCount, nil
}

func (s *PostgresStorage) CreateBookmark(ctx context.Context, title, url string) (*model.Bookmark, error) {
	var bm model.Bookmark

	row := s.db.QueryRowContext(ctx, "INSERT INTO bookmarks (title, url) VALUES($1, $2) RETURNING id, url, title, created_at, updated_at", title, url)

	err := row.Scan(&bm.ID, &bm.URL, &bm.Title, &bm.CreatedAt, &bm.UpdatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolation {
			return nil, storage.ErrExists
		}

		return nil, fmt.Errorf("failed to create bookmark: %w", err)
	}

	return &bm, nil
}

func (s *PostgresStorage) EditBookmark(ctx context.Context, id int, title, url string) (*model.Bookmark, error) {
	var bm model.Bookmark

	row := s.db.QueryRowContext(ctx, "UPDATE bookmarks SET title=$1, url=$2, updated_at=NOW() WHERE id=$3 RETURNING id, url, title, created_at, updated_at", title, url, id)

	err := row.Scan(&bm.ID, &bm.URL, &bm.Title, &bm.CreatedAt, &bm.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("failed to edit bookmark: %w", err)
	}

	return &bm, nil
}

func (s *PostgresStorage) DeleteBookmark(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM bookmarks WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get check deletion: %w", err)
	}

	if rowAffected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *PostgresStorage) BookmarkExist(ctx context.Context, url string) (int, bool, error) {
	var id int
	var found bool

	err := s.db.QueryRowContext(ctx, `SELECT id, true FROM bookmarks WHERE url = $1 LIMIT 1`, url).Scan(&id, &found)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}

		return 0, false, fmt.Errorf("failed to find bookmark: %w", err)
	}

	return id, found, nil
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	return nil
}

func (s *PostgresStorage) Close() error {
	slog.Info("closing database connection")

	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
