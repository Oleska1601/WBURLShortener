package models

import "time"

type URL struct {
	ID          int64     `json:"id"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}
