package postgres

import (
	"context"
	"strings"

	"gorm.io/gorm/clause"
)

type PlayerRosterStats struct {
	RosterSize  int64 `gorm:"column:roster_size"`
	Active7d    int64 `gorm:"column:active_7d"`
	Active30d   int64 `gorm:"column:active_30d"`
	Inactive30d int64 `gorm:"column:inactive_30d"`
}

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

func (s *Postgres) SearchGuilds(ctx context.Context, region Region, prefix string) ([]string, error) {
	var guilds []string
	err := s.db.WithContext(ctx).
		Model(&PlayerStatsLatest{}).
		Distinct("guild_name").
		Where("region = ? AND LOWER(guild_name) LIKE ?", region, strings.ToLower(prefix)+"%").
		Order("guild_name ASC").
		Limit(6).
		Pluck("guild_name", &guilds).Error
	return guilds, err
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

func (s *Postgres) GetAllianceRosterStats(ctx context.Context, region Region, allianceName string) (*PlayerRosterStats, error) {
	var stats PlayerRosterStats
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) AS roster_size,
			COUNT(*) FILTER (WHERE GREATEST(killboard_last_activity, other_last_activity) > NOW() - INTERVAL '7 days') AS active_7d,
			COUNT(*) FILTER (WHERE GREATEST(killboard_last_activity, other_last_activity) > NOW() - INTERVAL '30 days') AS active_30d,
			COUNT(*) FILTER (WHERE GREATEST(killboard_last_activity, other_last_activity) < NOW() - INTERVAL '30 days') AS inactive_30d
		FROM player_stats_latest
		WHERE region = ?
			AND alliance_name = ?
	`, region, allianceName).Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (s *Postgres) GetGuildRosterStats(ctx context.Context, region Region, guildName string) (*PlayerRosterStats, error) {
	var stats PlayerRosterStats
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) AS roster_size,
			COUNT(*) FILTER (WHERE GREATEST(killboard_last_activity, other_last_activity) > NOW() - INTERVAL '7 days') AS active_7d,
			COUNT(*) FILTER (WHERE GREATEST(killboard_last_activity, other_last_activity) > NOW() - INTERVAL '30 days') AS active_30d,
			COUNT(*) FILTER (WHERE GREATEST(killboard_last_activity, other_last_activity) < NOW() - INTERVAL '30 days') AS inactive_30d
		FROM player_stats_latest
		WHERE region = ?
			AND guild_name = ?
	`, region, guildName).Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
