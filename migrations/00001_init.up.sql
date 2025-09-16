CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE bookmarks (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_bookmarks_title_trgm
ON bookmarks USING GIN (title gin_trgm_ops);
