package postgres

import (
	"context"

	"gorm.io/gorm/clause"
)

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

func (s *Postgres) UpsertPlayerStatsLatest(stats []PlayerStatsLatest) error {
	if len(stats) == 0 {
		return nil
	}

	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		UpdateAll: true,
	}).Create(&stats).Error
}
