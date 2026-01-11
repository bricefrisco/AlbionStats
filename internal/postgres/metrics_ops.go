package postgres

import "context"

func (s *Postgres) InsertMetrics(ctx context.Context) error {
	return s.db.WithContext(ctx).Exec(`
		INSERT INTO metrics (metric, ts, value)
		VALUES
			('players_total', now(), (SELECT COUNT(*) FROM player_stats_latest)),
			('snapshots', now(), (SELECT COUNT(*) FROM player_stats_snapshots))
	`).Error
}
