package controller

import (
	"context"

	"github.com/Oleska1601/WBURLShortener/internal/models"
)

type UsecaseInterface interface {
	URLUsecaseInterface
	AnalyticsUsecaseInterface
}

type URLUsecaseInterface interface {
	GetURL(context.Context, string) (string, error)
	CreateURL(context.Context, *models.URL) (*models.URL, error)
}

type AnalyticsUsecaseInterface interface {
	GetAnalytics(context.Context, string) (*models.AnAgregation, error)
	CreateAnalytics(context.Context, *models.Analytics) (int64, error)
}
