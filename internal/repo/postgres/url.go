package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Oleska1601/WBURLShortener/internal/errs"
	"github.com/Oleska1601/WBURLShortener/internal/models"
)

const (
	selectOriginalURL = `SELECT original_url FROM urls WHERE short_url = $1`
	insertURL         = `INSERT INTO urls (short_url, original_url) VALUES ($1, $2) RETURNING id`
)

func (r *PgRepo) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	row := r.db.QueryRowContext(ctx, selectOriginalURL, shortURL)
	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.NewNotFoundError(fmt.Sprintf("get original url by short_url %s", shortURL))
		}

		return "", fmt.Errorf("row scan: %w", err)
	}

	return originalURL, nil
}

func (r *PgRepo) CreateURL(ctx context.Context, url *models.URL) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, insertURL,
		url.ShortURL,
		url.OriginalURL).
		Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, errs.NewAlreadyExistsError(fmt.Sprintf("create url with short_url %s", url.ShortURL))
		}

		return 0, fmt.Errorf("create url: %w", err)
	}

	return id, nil
}
