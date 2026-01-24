package postgres

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TopAllianceStats struct {
	AllianceName   string `gorm:"column:alliance_name"`
	TotalKillFame  int64  `gorm:"column:total_kill_fame"`
	TotalDeathFame int64  `gorm:"column:total_death_fame"`
	TotalKills     int64  `gorm:"column:total_kills"`
	TotalDeaths    int64  `gorm:"column:total_deaths"`
}

func (p *Postgres) InsertBattleAllianceStats(stats []BattleAllianceStats) error {
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

func (p *Postgres) UpdateBattleAllianceStats(stats []BattleAllianceStats) error {
	if len(stats) == 0 {
		return nil
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			updates := make(map[string]interface{})
			updates["death_fame"] = stat.DeathFame
			updates["ip"] = stat.IP

			if err := tx.Model(&BattleAllianceStats{}).
				Where("region = ? AND battle_id = ? AND alliance_name = ?", stat.Region, stat.BattleID, stat.AllianceName).
				Updates(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Postgres) GetBattleSummariesByAlliance(region string, allianceName string, playerCount int, limit int, offset int) ([]BattleSummary, error) {
	var summaries []BattleSummary
	err := p.db.Raw(`
SELECT bs.*
FROM (
    SELECT region, battle_id, start_time
    FROM battle_alliance_stats
    WHERE region = ?
      AND alliance_name = ?
      AND player_count >= ?
    ORDER BY start_time DESC
    LIMIT ? OFFSET ?
) bas
JOIN battle_summary bs
  ON bs.region = bas.region
 AND bs.battle_id = bas.battle_id
ORDER BY bas.start_time DESC;
	`, region, allianceName, playerCount, limit, offset).Scan(&summaries).Error

	return summaries, err
}

func (p *Postgres) GetBattleAllianceStatsByIDs(ctx context.Context, region string, battleIDs []int64) ([]BattleAllianceStats, error) {
	var stats []BattleAllianceStats
	err := p.db.WithContext(ctx).Where("region = ? AND battle_id IN ?", region, battleIDs).
		Order("player_count DESC").
		Find(&stats).Error

	return stats, err
}

func (p *Postgres) GetTopAlliances(region string, limit int, offset int) ([]TopAllianceStats, error) {
	var stats []TopAllianceStats
	err := p.db.Raw(`
		SELECT
			alliance_name,
			SUM(kill_fame) AS total_kill_fame,
			SUM(COALESCE(death_fame, 0)) AS total_death_fame,
			SUM(kills) AS total_kills,
			SUM(deaths) AS total_deaths
		FROM battle_guild_stats
		WHERE region = ?
			AND start_time >= NOW() - INTERVAL '30 days'
			AND alliance_name <> ''
		GROUP BY alliance_name
		ORDER BY total_kill_fame DESC
		LIMIT ? OFFSET ?
	`, region, limit, offset).Scan(&stats).Error

	return stats, err
}
