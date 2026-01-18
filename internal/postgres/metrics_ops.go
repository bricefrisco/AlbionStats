package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (s *Postgres) InsertPlayersTotalAndSnapshotMetrics(ctx context.Context) error {
	return s.db.WithContext(ctx).Exec(`
		INSERT INTO metrics (metric, ts, value)
		VALUES
			('players_total', now(), (SELECT COUNT(*) FROM player_stats_latest)),
			('snapshots', now(), (SELECT COUNT(*) FROM player_stats_snapshots))
	`).Error
}

func (s *Postgres) InsertActivePlayersMetrics(ctx context.Context) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Insert regional active player counts
		if err := tx.Exec(`
			INSERT INTO metrics (metric, ts, value)
			SELECT
				'active_players_24h_' || region as metric,
				now() as ts,
				COUNT(DISTINCT player_id) AS value
			FROM player_stats_snapshots
			WHERE ts >= now() - interval '24 hours'
				AND (
					killboard_last_activity >= now() - interval '24 hours'
					OR other_last_activity >= now() - interval '24 hours'
				)
			GROUP BY region
		`).Error; err != nil {
			return err
		}

		// Insert total active players count
		return tx.Exec(`
			INSERT INTO metrics (metric, ts, value)
			VALUES
				('active_players_24h', now(), (
					SELECT COUNT(DISTINCT player_id)
					FROM player_stats_snapshots
					WHERE ts >= now() - interval '24 hours'
						AND (
							killboard_last_activity >= now() - interval '24 hours'
							OR other_last_activity >= now() - interval '24 hours'
						)
				))
		`).Error
	})
}

func (s *Postgres) GetMetrics(ctx context.Context, metricId, granularity string) ([]int64, []int64, error) {
	var query string
	var args []interface{}

	switch granularity {
	case "1w":
		query = `
			SELECT
				time_bucket('6 hours', ts) AS timestamp,
				max(value) AS value
			FROM metrics
			WHERE metric = $1 AND ts >= NOW() - INTERVAL '1 week'
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	case "1m":
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
