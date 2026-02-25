package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Oleska1601/WBURLShortener/internal/errs"
	"github.com/Oleska1601/WBURLShortener/internal/models"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func (s *Service) getOriginalURL(ctx context.Context, shortURL string) (string, error) {
	cachedURL, err := s.cache.GetValue(ctx, shortURL)
	if err != nil && !errors.Is(err, redis.NoMatches) {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "getShortURL s.cache.GetValue").
			Str("short_url", shortURL).
			Msg("failed to get cache value")
	}

	if cachedURL != "" {
		return cachedURL, nil
	}

	originalURL, err := s.repo.GetOriginalURL(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("get url: %w", err)
	}

	return originalURL, nil
}

func (s *Service) GetURL(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := s.getOriginalURL(ctx, shortURL)
	if err != nil {
		return "", err
	}

	if err := s.cache.SetValue(ctx, shortURL, originalURL); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "GetURL s.cache.SetValue").
			Str("short_url", shortURL).
			Msg("failed to set cache value")
	}

	return originalURL, nil
}

func (s *Service) createURLWithInput(ctx context.Context, input *models.URL) (*models.URL, error) {
	_, err := s.getOriginalURL(ctx, input.ShortURL)
	if err == nil {
		return nil, errs.NewAlreadyExistsError(fmt.Sprintf("create url with short_url %s", input.ShortURL))
	}

	if !errors.Is(err, errs.NotFoundError) {
		return nil, fmt.Errorf("check url exists: %w", err)
	}

	id, err := s.repo.CreateURL(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create url: %w", err)
	}

	if err := s.cache.SetValue(ctx, input.ShortURL, input.OriginalURL); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "CreateShortURL s.cache.SetValue").
			Str("short_url", input.ShortURL).
			Msg("failed to set cache value")
	}

	input.ID = id
	return input, nil
}

func (s *Service) CreateURL(ctx context.Context, input *models.URL) (*models.URL, error) {
	// проверка, ввел ли пользователь ссылку
	if input.ShortURL != "" {
		return s.createURLWithInput(ctx, input)
	}

	// макс количество попыток генерации короткой ссылки
	// при коллизии (ссылка уже существует) генерируем новую,
	// но не более maxAttempts раз (чтобы избежать бесконечной генерации и зависания запроса)
	const maxAttemts = 5
	for range maxAttemts {
		input.ShortURL = generateShortURL()
		result, err := s.createURLWithInput(ctx, input)
		if err != nil {
			if errors.Is(err, errs.AlreadyExistsError) {
				continue
			}

			return nil, err
		}

		return result, nil
	}

	return nil, errs.NewConflictError("failed to generate unique short_url")
}
