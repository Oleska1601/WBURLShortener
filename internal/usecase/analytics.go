package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Oleska1601/WBURLShortener/internal/models"
	"github.com/go-redis/redis"
	"github.com/wb-go/wbf/zlog"
)

func (u *Usecase) GetAnalytics(ctx context.Context, shortURL string) (*models.AnAgregation, error) {
	cachedAnalytics, err := u.cache.GetValue(ctx, "analytics:"+shortURL)
	if err != nil && !errors.Is(err, redis.Nil) {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "GetAnalytics u.cache.GetValue").
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
	totalCount, err := u.repo.GetAnTotalCount(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("get analytics total count: %w", err)
	}

	dayCount, err := u.repo.GetAnDayCount(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("get analytics day count: %w", err)
	}

	monthCount, err := u.repo.GetAnMonthCount(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("get analytics month count: %w", err)
	}

	userAgentCount, err := u.repo.GetAnUserAgentCount(ctx, shortURL)
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
		if err := u.cache.SetValue(ctx, "analytics:"+shortURL, string(byteAn)); err != nil {
			zlog.Logger.Warn().
				Err(err).
				Str("path", "GetAnalytics u.cache.SetValue").
				Str("short_url", shortURL).
				Msg("failed to set cache value")
		}
	}

	return analytics, nil
}

func (u *Usecase) CreateAnalytics(ctx context.Context, analytics *models.Analytics) (int64, error) {
	id, err := u.repo.CreateAn(ctx, analytics)
	if err != nil {
		return 0, fmt.Errorf("create analytics: %w", err)
	}

	if err := u.cache.DeleteValue(ctx, "analytics:"+analytics.ShortURL); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "CreateAnalytics u.cache.DeleteValue").
			Str("short_url", analytics.ShortURL).
			Msg("failed to delete cache value")
	}

	return id, nil
}
