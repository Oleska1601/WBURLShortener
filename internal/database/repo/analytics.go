package repo

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
	fmt.Println(shortURL)
	row := r.db.Master.QueryRowContext(ctx, selectAnCount, shortURL)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}

	return count, nil
}

func (r *PgRepo) GetAnDayCount(ctx context.Context, shortURL string) (map[string]int, error) {
	rows, err := r.db.Master.QueryContext(ctx, selectAnDayCount, shortURL)
	if err != nil {
		return nil, fmt.Errorf("r.db.Master.QueryContext: %w", err)
	}

	dayCount := make(map[string]int)
	var date string
	var count int
	for rows.Next() {
		if err := rows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		dayCount[date] = count
	}

	return dayCount, nil
}

func (r *PgRepo) GetAnMonthCount(ctx context.Context, shortURL string) (map[string]int, error) {
	rows, err := r.db.Master.QueryContext(ctx, selectAnMonthCount, shortURL)
	if err != nil {
		return nil, fmt.Errorf("r.db.Master.QueryContext: %w", err)
	}

	monthCount := make(map[string]int)
	var date string
	var count int
	for rows.Next() {
		if err := rows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		monthCount[date] = count
	}

	return monthCount, nil
}

func (r *PgRepo) GetAnUserAgentCount(ctx context.Context, shortURL string) (map[string]int, error) {
	rows, err := r.db.Master.QueryContext(ctx, selectAnUserAgentCount, shortURL)
	if err != nil {
		return nil, fmt.Errorf("r.db.Master.QueryContext: %w", err)
	}
	userAgentCount := make(map[string]int)
	var date string
	var count int
	for rows.Next() {
		if err := rows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		userAgentCount[date] = count
	}

	return userAgentCount, nil
}

func (r *PgRepo) CreateAn(ctx context.Context, analytics *models.Analytics) (int64, error) {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("r.db.Master.BeginTx: %w", err)
	}

	defer r.rollbackTransaction(tx)

	var id int64
	err = tx.QueryRowContext(ctx, insertAn,
		analytics.ShortURL, analytics.UserAgent, analytics.IP).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("tx.QueryRowContext Scan: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("tx.Commit: %w", err)
	}

	return id, nil
}
