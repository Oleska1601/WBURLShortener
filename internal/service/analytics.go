package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Oleska1601/WBURLShortener/internal/models"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func (s *Service) GetAnalytics(ctx context.Context, shortURL string) (*models.AnAgregation, error) {
	cachedAnalytics, err := s.cache.GetValue(ctx, "analytics:"+shortURL)
	if err != nil && !errors.Is(err, redis.NoMatches) {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "GetAnalytics s.cache.GetValue").
			Str("short_url", shortURL).
			Msg("failed to get cache value")
	}

	if cachedAnalytics != "" {
		var analytics *models.AnAgregation
		if err := json.Unmarshal([]byte(cachedAnalytics), &analytics); err != nil {
			zlog.Logger.Warn().
				Err(err).
				Str("path", "GetAnalytics json.Unmarshal").
				Str("short_url", shortURL).
				Msg("failed to unmarshal json")
			// если unmarshal не удался -> собираем аналитику заново
		} else {
			return analytics, nil
		}
	}
	totalCount, err := s.repo.GetAnTotalCount(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("get analytics total count: %w", err)
	}

	dayCount, err := s.repo.GetAnDayCount(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("get analytics day count: %w", err)
	}

	monthCount, err := s.repo.GetAnMonthCount(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("get analytics month count: %w", err)
	}

	userAgentCount, err := s.repo.GetAnUserAgentCount(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("get analytics user agent count: %w", err)
	}

	analytics := &models.AnAgregation{
		ShortURL:       shortURL,
		TotalCount:     totalCount,
		DayCount:       dayCount,
		MonthCount:     monthCount,
		UserAgentCount: userAgentCount,
	}

	// связано с ошибками кеша, но результат успешно получен -> warn, а пользователю все равно возвращаем значение
	if byteAn, err := json.Marshal(analytics); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "GetAnalytics json.Marshal").
			Str("short_url", shortURL).
			Msg("failed to marshal json")
	} else {
		if err := s.cache.SetValue(ctx, "analytics:"+shortURL, string(byteAn)); err != nil {
			zlog.Logger.Warn().
				Err(err).
				Str("path", "GetAnalytics s.cache.SetValue").
				Str("short_url", shortURL).
				Msg("failed to set cache value")
		}
	}

	return analytics, nil
}

func (s *Service) CreateAnalytics(ctx context.Context, analytics *models.Analytics) (int, error) {
	id, err := s.repo.CreateAn(ctx, analytics)
	if err != nil {
		return 0, fmt.Errorf("create analytics: %w", err)
	}

	// удаление текущего значения аналитики, поскольку добавилось новое знач
	if err := s.cache.DeleteValue(ctx, "analytics:"+analytics.ShortURL); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "CreateAnalytics s.cache.DeleteValue").
			Str("short_url", analytics.ShortURL).
			Msg("failed to delete cache value")
	}

	return id, nil
}
