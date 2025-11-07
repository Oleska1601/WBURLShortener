package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Oleska1601/WBURLShortener/internal/apperrors"
	"github.com/Oleska1601/WBURLShortener/internal/models"
	"github.com/go-redis/redis"
	"github.com/wb-go/wbf/zlog"
)

// get original url
func (u *Usecase) GetURL(ctx context.Context, shortURL string) (string, error) {
	cachedURL, err := u.cache.GetValue(ctx, shortURL)
	if err != nil && !errors.Is(err, redis.Nil) {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "getShortURL u.cache.GetValue").
			Str("short_url", shortURL).
			Msg("failed to get cache value")
	}

	if cachedURL != "" {
		return cachedURL, nil
	}

	originalURL, err := u.repo.GetOriginalURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apperrors.NewNotFoundError("url not found")
		}

		return "", fmt.Errorf("get url: %w", err)
	}

	if err := u.cache.SetValue(ctx, shortURL, originalURL); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "GetURL u.cache.SetValue").
			Str("short_url", shortURL).
			Msg("failed to set cache value")
	}

	return originalURL, nil

}

func (u *Usecase) CreateURL(ctx context.Context, input *models.URL) (*models.URL, error) {
	//check if input short url already exists
	if input.ShortURL == "" {
		input.ShortURL = generateShortURL()
	} else {
		shortURL := input.ShortURL
		_, err := u.cache.GetValue(ctx, shortURL)
		if err == nil {
			return nil, apperrors.NewAlreadyExistsError("input short url already exists")
		}

		if !errors.Is(err, redis.Nil) {
			zlog.Logger.Warn().
				Err(err).
				Str("path", "CreateShortURL u.cache.GetValue").
				Str("short_url", shortURL).
				Msg("failed to get cache value")
		}

		_, err = u.repo.GetOriginalURL(ctx, shortURL)
		if err == nil {
			return nil, apperrors.NewAlreadyExistsError("input short url already exists")
		}

		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get url: %w", err)
		}

	}

	id, err := u.repo.CreateURL(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create url: %w", err)
	}

	if err := u.cache.SetValue(ctx, input.ShortURL, input.OriginalURL); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "CreateShortURL u.cache.SetValue").
			Str("short_url", input.ShortURL).
			Msg("failed to set cache value")
	}

	return &models.URL{
		ID:       id,
		ShortURL: input.ShortURL,
	}, nil
}
