package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
		    SELECT region, battle_id
		    FROM battle_alliance_stats
		    WHERE region = ?
		    AND alliance_name = ?
		    AND player_count >= ?
		) bas
		JOIN battle_summary bs
		ON bs.region = bas.region
		AND bs.battle_id = bas.battle_id
		ORDER BY bs.start_time DESC
		LIMIT ? OFFSET ?
	`, region, allianceName, playerCount, limit, offset).Scan(&summaries).Error

	return summaries, err
}

func (p *Postgres) GetBattleAllianceStatsByIDs(region string, battleIDs []int64) ([]BattleAllianceStats, error) {
	var stats []BattleAllianceStats
	err := p.db.Where("region = ? AND battle_id IN ?", region, battleIDs).
		Order("player_count DESC").
		Find(&stats).Error

	return stats, err
}
