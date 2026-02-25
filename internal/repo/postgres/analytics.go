package postgres

import (
	"context"
	"fmt"

	"github.com/Oleska1601/WBURLShortener/internal/models"
)

const (
	insertAn         = `INSERT INTO analytics (short_url, user_agent, ip) VALUES ($1, $2, $3) RETURNING id`
	selectAnCount    = `SELECT count(*) AS count FROM analytics WHERE short_url = $1`
	selectAnDayCount = `SELECT TO_CHAR("requested_at", 'YYYY-MM-DD') AS day, count(*) AS count FROM analytics 
						WHERE short_url = $1 
						GROUP BY day`
	selectAnMonthCount = `SELECT TO_CHAR("requested_at", 'YYYY-MM') AS month, count(*) AS count FROM analytics 
						WHERE short_url = $1
						GROUP BY month`
	selectAnUserAgentCount = `SELECT user_agent, count(*) AS count FROM analytics 
							WHERE short_url = $1
							GROUP BY user_agent`
)

func (r *PgRepo) GetAnTotalCount(ctx context.Context, shortURL string) (int, error) {
	row := r.db.QueryRowContext(ctx, selectAnCount, shortURL)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("row scan: %w", err)
	}

	return count, nil
}

func (r *PgRepo) GetAnDayCount(ctx context.Context, shortURL string) (map[string]int, error) {
	rows, err := r.db.QueryContext(ctx, selectAnDayCount, shortURL)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}

	defer rows.Close()

	dayCount := make(map[string]int)
	var date string
	var count int
	for rows.Next() {
		err := rows.Scan(&date, &count)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		dayCount[date] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return dayCount, nil
}

func (r *PgRepo) GetAnMonthCount(ctx context.Context, shortURL string) (map[string]int, error) {
	rows, err := r.db.QueryContext(ctx, selectAnMonthCount, shortURL)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}

	defer rows.Close()

	monthCount := make(map[string]int)
	var date string
	var count int
	for rows.Next() {
		err := rows.Scan(&date, &count)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		monthCount[date] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return monthCount, nil
}

func (r *PgRepo) GetAnUserAgentCount(ctx context.Context, shortURL string) (map[string]int, error) {
	rows, err := r.db.QueryContext(ctx, selectAnUserAgentCount, shortURL)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}

	defer rows.Close()

	userAgentCount := make(map[string]int)
	var userAgent string
	var count int
	for rows.Next() {
		err := rows.Scan(&userAgent, &count)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		userAgentCount[userAgent] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return userAgentCount, nil
}

func (r *PgRepo) CreateAn(ctx context.Context, analytics *models.Analytics) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, insertAn,
		analytics.ShortURL,
		analytics.UserAgent,
		analytics.IP).
		Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create an: %w", err)
	}

	return id, nil
}
