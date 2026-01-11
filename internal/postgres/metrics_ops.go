package postgres

import (
	"context"
)

func (s *Postgres) InsertMetrics(ctx context.Context) error {
	return s.db.WithContext(ctx).Exec(`
		INSERT INTO metrics (metric, ts, value)
		VALUES
			('players_total', now(), (SELECT COUNT(*) FROM player_stats_latest)),
			('snapshots', now(), (SELECT COUNT(*) FROM player_stats_snapshots))
	`).Error
}

func (s *Postgres) SearchPlayers(ctx context.Context, region Region, prefix string, limit int) ([]PlayerStatsLatest, error) {
	var players []PlayerStatsLatest
	err := s.db.WithContext(ctx).
		Select("player_id", "name", "guild_name", "alliance_name").
		Where("region = ? AND LOWER(name) LIKE LOWER(?)", region, prefix+"%").
		Limit(limit).
		Order("name ASC").
		Find(&players).Error
	return players, err
}
