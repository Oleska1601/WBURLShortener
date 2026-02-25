package service

import (
	"context"

	"github.com/Oleska1601/WBURLShortener/internal/models"
)

type CacheI interface {
	GetValue(context.Context, string) (string, error)
	SetValue(context.Context, string, any) error
	DeleteValue(context.Context, string) error
}

type RepoI interface {
	RepoURLI
	RepoAnI
}

type RepoURLI interface {
	GetOriginalURL(context.Context, string) (string, error)
	CreateURL(context.Context, *models.URL) (int, error)
}

type RepoAnI interface {
	GetAnTotalCount(context.Context, string) (int, error)
	GetAnDayCount(context.Context, string) (map[string]int, error)
	GetAnMonthCount(context.Context, string) (map[string]int, error)
	GetAnUserAgentCount(context.Context, string) (map[string]int, error)

	CreateAn(context.Context, *models.Analytics) (int, error)
}
