package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TopGuildStats struct {
	GuildName      string `gorm:"column:guild_name"`
	TotalKillFame  int64  `gorm:"column:total_kill_fame"`
	TotalDeathFame int64  `gorm:"column:total_death_fame"`
	TotalKills     int64  `gorm:"column:total_kills"`
	TotalDeaths    int64  `gorm:"column:total_deaths"`
}

type GuildBattleSummary struct {
	Battles        int64     `gorm:"column:battles"`
	TotalKills     int64     `gorm:"column:total_kills"`
	TotalDeaths    int64     `gorm:"column:total_deaths"`
	TotalKillFame  int64     `gorm:"column:total_kill_fame"`
	TotalDeathFame int64     `gorm:"column:total_death_fame"`
	MaxPlayers     int32     `gorm:"column:max_players"`
	LastBattleAt   time.Time `gorm:"column:last_battle_at"`
}

func (p *Postgres) InsertBattleGuildStats(stats []BattleGuildStats) error {
	if len(stats) == 0 {
		return nil
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&stat).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Postgres) UpdateBattleGuildStats(stats []BattleGuildStats) error {
	if len(stats) == 0 {
		return nil
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			updates := make(map[string]interface{})
			updates["death_fame"] = stat.DeathFame
			updates["ip"] = stat.IP

			if err := tx.Model(&BattleGuildStats{}).
				Where("region = ? AND battle_id = ? AND guild_name = ?", stat.Region, stat.BattleID, stat.GuildName).
				Updates(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Postgres) GetBattleSummariesByGuild(region string, guildName string, playerCount int, limit int, offset int) ([]BattleSummary, error) {
	var summaries []BattleSummary
	err := p.db.Raw(`
SELECT bs.*
FROM (
    SELECT region, battle_id, start_time
    FROM battle_guild_stats
    WHERE region = ?
      AND guild_name = ?
      AND player_count >= ?
    ORDER BY start_time DESC
    LIMIT ? OFFSET ?
) bgs
JOIN battle_summary bs
  ON bs.region = bgs.region
 AND bs.battle_id = bgs.battle_id
ORDER BY bgs.start_time DESC;
	`, region, guildName, playerCount, limit, offset).Scan(&summaries).Error

	return summaries, err
}

func (p *Postgres) GetBattleGuildStatsByIDs(ctx context.Context, region string, battleIDs []int64) ([]BattleGuildStats, error) {
	var stats []BattleGuildStats
	err := p.db.WithContext(ctx).Where("region = ? AND battle_id IN ?", region, battleIDs).
		Order("player_count DESC").
		Find(&stats).Error

	return stats, err
}

func (p *Postgres) GetTopGuilds(region string, limit int, offset int) ([]TopGuildStats, error) {
	var stats []TopGuildStats
	err := p.db.Raw(`
		SELECT
			guild_name,
			SUM(kill_fame) AS total_kill_fame,
			SUM(COALESCE(death_fame, 0)) AS total_death_fame,
			SUM(kills) AS total_kills,
			SUM(deaths) AS total_deaths
		FROM battle_guild_stats
		WHERE region = ?
			AND start_time >= NOW() - INTERVAL '30 days'
		GROUP BY guild_name
		ORDER BY total_kill_fame DESC
		LIMIT ? OFFSET ?
	`, region, limit, offset).Scan(&stats).Error

	return stats, err
}

func (p *Postgres) GetGuildBattleSummary(ctx context.Context, region string, guildName string) (*GuildBattleSummary, error) {
	var summary GuildBattleSummary
	err := p.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(DISTINCT battle_id) AS battles,
			SUM(kills) AS total_kills,
			SUM(deaths) AS total_deaths,
			SUM(kill_fame) AS total_kill_fame,
			SUM(death_fame) AS total_death_fame,
			MAX(player_count) AS max_players,
			MAX(start_time) AS last_battle_at
		FROM battle_guild_stats
		WHERE region = ?
			AND guild_name = ?
	`, region, guildName).Scan(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}
