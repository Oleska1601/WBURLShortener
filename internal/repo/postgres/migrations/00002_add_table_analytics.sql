-- +goose Up
CREATE TABLE IF NOT EXISTS analytics (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(6) NOT NULL REFERENCES urls(short_url) ON DELETE CASCADE,
    user_agent VARCHAR(255) NOT NULL,
    IP VARCHAR(255) NOT NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS analytics;