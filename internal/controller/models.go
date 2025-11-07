package controller

import "github.com/Oleska1601/WBURLShortener/internal/models"

type CreateShortURLResponse struct {
	ID       int64  `json:"id"`
	ShortURL string `json:"short_url"`
}

func toCreateShortURLResponse(url *models.URL) *CreateShortURLResponse {
	return &CreateShortURLResponse{
		ID:       url.ID,
		ShortURL: url.ShortURL,
	}
}

type GetAnalyticsResponse struct {
	ShortURL       string         `json:"short_url"`
	TotalCount     int            `json:"total_count"`
	DayCount       map[string]int `json:"day_count"`
	MonthCount     map[string]int `json:"month_count"`
	UserAgentCount map[string]int `json:"user_agent_count"`
}

func toAnalyticsResponse(analytics *models.AnAgregation) *GetAnalyticsResponse {
	return &GetAnalyticsResponse{
		ShortURL:       analytics.ShortURL,
		TotalCount:     analytics.TotalCount,
		DayCount:       analytics.DayCount,
		MonthCount:     analytics.MonthCount,
		UserAgentCount: analytics.UserAgentCount,
	}
}
