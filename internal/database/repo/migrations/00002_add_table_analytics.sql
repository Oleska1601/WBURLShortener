-- +goose Up
CREATE TABLE IF NOT EXISTS analytics (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(7) NOT NULL REFERENCES urls(short_url) ON DELETE CASCADE, 
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_agent VARCHAR(255) NOT NULL,
    IP VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS analytics;