-- +goose Up
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code TEXT NOT NULL,
    original_url TEXT NOT NULL,
    created_by TEXT,
    expired_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE urls;