-- +goose Up
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(6) UNIQUE NOT NULL,
    original_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS urls;