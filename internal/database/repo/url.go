package repo

import (
	"context"
	"fmt"

	"github.com/Oleska1601/WBURLShortener/internal/models"
)

const (
	getOriginalURL = `SELECT original_url FROM urls WHERE short_url = $1`
	insertURL      = `INSERT INTO urls (short_url, original_url) VALUES ($1, $2) RETURNING id`
)

func (r *PgRepo) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	row := r.db.Master.QueryRowContext(ctx, getOriginalURL, shortURL)
	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		return "", fmt.Errorf("row.Scan: %w", err)
	}

	return originalURL, nil
}

func (r *PgRepo) CreateURL(ctx context.Context, url *models.URL) (int64, error) {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("r.db.Master.BeginTx: %w", err)
	}

	defer r.rollbackTransaction(tx)

	var id int64
	err = tx.QueryRowContext(ctx, insertURL,
		url.ShortURL, url.OriginalURL).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("tx.QueryRowContext Scan: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("tx.Commit: %w", err)
	}

	return id, nil
}
