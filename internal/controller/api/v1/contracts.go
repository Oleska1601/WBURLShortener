package v1

import (
	"context"

	"github.com/Oleska1601/WBURLShortener/internal/models"
)

type ServiceI interface {
	URLI
	AnalyticsI
}

type URLI interface {
	GetURL(context.Context, string) (string, error)
	CreateURL(context.Context, *models.URL) (*models.URL, error)
}

type AnalyticsI interface {
	GetAnalytics(context.Context, string) (*models.AnAgregation, error)
	CreateAnalytics(context.Context, *models.Analytics) (int, error)
}
