package models

import "time"

type Analytics struct {
	ShortURL    string    `json:"short_url"`
	RequestedAt time.Time `json:"requested_at"`
	UserAgent   string    `json:"user_agent"`
	IP          string    `json:"ip"`
}

// analytics agregation
type AnAgregation struct {
	ShortURL       string         `json:"short_url"`
	TotalCount     int            `json:"total_count"`
	DayCount       map[string]int `json:"day_count"`
	MonthCount     map[string]int `json:"month_count"`
	UserAgentCount map[string]int `json:"user_agent_count"`
}
