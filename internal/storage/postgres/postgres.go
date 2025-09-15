package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/haadi-coder/bookmark-manager/internal/model"
	"github.com/haadi-coder/bookmark-manager/internal/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func New(path url.URL) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", path.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(3)

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) GetBookmarks(ctx context.Context, limit, offset int, search string) ([]*model.Bookmark, int, error) {
	var rows *sql.Rows

	stmt := sq.
		Select("*", "COUNT(*) OVER() AS total_count").
		From("bookmarks").
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		RunWith(s.db)

	if search != "" {
		stmt = stmt.Where(
			sq.Or{
				sq.ILike{"url": "%" + search + "%"},
				sq.ILike{"title": "%" + search + "%"},
			})
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bookmarks rows: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var totalCount int
	bookmarks := []*model.Bookmark{}
	for rows.Next() {
		var bm model.Bookmark

		if err := rows.Scan(&bm.ID, &bm.URL, &bm.Title, &bm.CreatedAt, &bm.UpdatedAt, &totalCount); err != nil {
			return nil, 0, fmt.Errorf("failed to scan bookmark: %w", err)
		}

		bookmarks = append(bookmarks, &bm)
	}

	return bookmarks, totalCount, nil
}

func (s *PostgresStorage) CreateBookmark(ctx context.Context, title, url string) (*model.Bookmark, error) {
	const uniqueViolation = "23505"

	stmt := sq.
		Insert("bookmarks").
		Columns("title", "url").
		Values(title, url).
		Suffix("RETURNING id, url, title, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		RunWith(s.db)

	row := stmt.QueryRowContext(ctx)

	var bm model.Bookmark
	if err := row.Scan(&bm.ID, &bm.URL, &bm.Title, &bm.CreatedAt, &bm.UpdatedAt); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolation {
			return nil, storage.ErrExists
		}

		return nil, fmt.Errorf("failed to create bookmark: %w", err)
	}

	return &bm, nil
}

func (s *PostgresStorage) EditBookmark(ctx context.Context, id int, title, url string) (*model.Bookmark, error) {
	stmt := sq.
		Update("bookmarks").
		Set("title", title).
		Set("url", url).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, url, title, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		RunWith(s.db)

	row := stmt.QueryRowContext(ctx)

	var bm model.Bookmark
	if err := row.Scan(&bm.ID, &bm.URL, &bm.Title, &bm.CreatedAt, &bm.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("failed to edit bookmark: %w", err)
	}

	return &bm, nil
}

func (s *PostgresStorage) DeleteBookmark(ctx context.Context, id int) error {
	stmt := sq.
		Delete("bookmarks").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(s.db)

	result, err := stmt.ExecContext(ctx)
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

	stmt := sq.
		Select("id", "true").
		From("bookmarks").
		Where(sq.Eq{"url": url}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		RunWith(s.db)

	err := stmt.QueryRowContext(ctx).Scan(&id, &found)
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
