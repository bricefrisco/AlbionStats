package postgres

import (
	"context"
	"database/sql"
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
		return tx.Exec(`
			INSERT INTO metrics (metric, ts, value)
			SELECT
				'active_players_24h_' || region AS metric,
				now() AS ts,
				COUNT(*) AS value
			FROM player_stats_latest
			WHERE last_activity >= now() - interval '24 hours'
			GROUP BY region
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

func (s *Postgres) GetDAUMetrics(ctx context.Context) ([]int64, []int64, []int64, []int64, error) {
	query := `
			SELECT
				time_bucket('6 hours', ts) AS timestamp,
				max(CASE WHEN metric = $1 THEN value END) AS americas,
				max(CASE WHEN metric = $2 THEN value END) AS europe,
				max(CASE WHEN metric = $3 THEN value END) AS asia
			FROM metrics
			WHERE metric IN ($1, $2, $3) AND ts >= NOW() - INTERVAL '1 week'
			GROUP BY 1
			ORDER BY 1`
	args := []interface{}{
		"active_players_24h_americas",
		"active_players_24h_europe",
		"active_players_24h_asia",
	}

	rows, err := s.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	var timestamps []int64
	var americas []int64
	var europe []int64
	var asia []int64
	for rows.Next() {
		var timestamp time.Time
		var americasValue sql.NullFloat64
		var europeValue sql.NullFloat64
		var asiaValue sql.NullFloat64
		if err := rows.Scan(&timestamp, &americasValue, &europeValue, &asiaValue); err != nil {
			return nil, nil, nil, nil, err
		}
		timestamps = append(timestamps, timestamp.UnixMilli())
		if americasValue.Valid {
			americas = append(americas, int64(americasValue.Float64))
		} else {
			americas = append(americas, 0)
		}
		if europeValue.Valid {
			europe = append(europe, int64(europeValue.Float64))
		} else {
			europe = append(europe, 0)
		}
		if asiaValue.Valid {
			asia = append(asia, int64(asiaValue.Float64))
		} else {
			asia = append(asia, 0)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, nil, nil, nil, err
	}

	return timestamps, americas, europe, asia, nil
}
