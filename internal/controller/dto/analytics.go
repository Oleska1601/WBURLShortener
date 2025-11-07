package dto

import "github.com/Oleska1601/WBURLShortener/internal/models"

func CreateAnalytics(shortURL, userAgent, ip string) *models.Analytics {
	return &models.Analytics{
		ShortURL:  shortURL,
		UserAgent: userAgent,
		IP:        ip,
	}
}
