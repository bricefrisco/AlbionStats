package postgres

import (
	"context"
	"strings"

	"gorm.io/gorm/clause"
)

func (s *Postgres) SearchPlayers(ctx context.Context, region Region, prefix string, limit int) ([]PlayerStatsLatest, error) {
	var players []PlayerStatsLatest
	err := s.db.WithContext(ctx).
		Select("player_id", "name", "guild_name", "alliance_name").
		Where("region = ? AND LOWER(name) LIKE ?", region, strings.ToLower(prefix)+"%").
		Limit(limit).
		Order("lower(name) ASC").
		Find(&players).Error
	return players, err
}

func (s *Postgres) SearchAlliances(ctx context.Context, region Region, prefix string) ([]string, error) {
	var alliances []string
	err := s.db.WithContext(ctx).
		Model(&PlayerStatsLatest{}).
		Distinct("alliance_name").
		Where("region = ? AND LOWER(alliance_name) LIKE ?", region, strings.ToLower(prefix)+"%").
		Order("alliance_name ASC").
		Limit(6).
		Pluck("alliance_name", &alliances).Error
	return alliances, err
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

func (s *Postgres) GetPlayerByName(ctx context.Context, region Region, name string) (*PlayerStatsLatest, error) {
	var player PlayerStatsLatest
	err := s.db.WithContext(ctx).
		Where("region = ? AND LOWER(name) = ?", region, strings.ToLower(name)).
		First(&player).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}
