package postgres

import (
	"context"
	"fmt"
	"time"
)

func (s *Postgres) InsertMetrics(ctx context.Context) error {
	return s.db.WithContext(ctx).Exec(`
		INSERT INTO metrics (metric, ts, value)
		VALUES
			('players_total', now(), (SELECT COUNT(*) FROM player_stats_latest)),
			('snapshots', now(), (SELECT COUNT(*) FROM player_stats_snapshots))
	`).Error
}

func (s *Postgres) GetMetrics(ctx context.Context, metricId, granularity string) ([]int64, []int64, error) {
	var query string
	var args []interface{}

	switch granularity {
	case "1w":
		// Aggregate to hourly buckets for 1 week
		query = `
			SELECT
				time_bucket('1 hour', ts) AS timestamp,
				max(value) AS value
			FROM metrics
			WHERE metric = $1 AND ts >= NOW() - INTERVAL '1 week'
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	case "1m":
		// Aggregate to daily buckets for 1 month
		query = `
			SELECT
				time_bucket('1 day', ts) AS timestamp,
				max(value) AS value
			FROM metrics
			WHERE metric = $1 AND ts >= NOW() - INTERVAL '1 month'
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	case "1y":
		// Aggregate to weekly buckets for 1 year
		query = `
			SELECT
				time_bucket('1 week', ts) AS timestamp,
				max(value) AS value
			FROM metrics
			WHERE metric = $1 AND ts >= NOW() - INTERVAL '1 year'
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	case "all":
		// Aggregate to monthly buckets for all time
		query = `
			SELECT
				time_bucket('1 month', ts) AS timestamp,
				max(value) AS value
			FROM metrics
			WHERE metric = $1
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	default:
		return nil, nil, fmt.Errorf("invalid granularity: %s", granularity)
	}

	rows, err := s.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var timestamps []int64
	var values []int64
	for rows.Next() {
		var timestamp time.Time
		var value float64
		if err := rows.Scan(&timestamp, &value); err != nil {
			return nil, nil, err
		}
		timestamps = append(timestamps, timestamp.UnixMilli())
		values = append(values, int64(value))
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return timestamps, values, nil
}
