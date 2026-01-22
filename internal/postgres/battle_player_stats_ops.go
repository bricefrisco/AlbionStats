package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p *Postgres) InsertBattlePlayerStats(stats []BattlePlayerStats) error {
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

func (p *Postgres) UpdateBattlePlayerStats(stats []BattlePlayerStats) error {
	if len(stats) == 0 {
		return nil
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			updates := make(map[string]interface{})
			updates["death_fame"] = stat.DeathFame
			updates["ip"] = stat.IP
			updates["weapon"] = stat.Weapon
			updates["damage"] = stat.Damage
			updates["heal"] = stat.Heal

			if err := tx.Model(&BattlePlayerStats{}).
				Where("region = ? AND battle_id = ? AND player_name = ?", stat.Region, stat.BattleID, stat.PlayerName).
				Updates(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Postgres) GetBattleSummariesByPlayer(region string, playerName string, playerCount int, limit int, offset int) ([]BattleSummary, error) {
	var summaries []BattleSummary
	err := p.db.Raw(`
		SELECT bs.*
		FROM (
		    SELECT region, battle_id
		    FROM battle_player_stats
		    WHERE region = ?
		    AND player_name = ?
		) bps
		JOIN battle_summary bs
		ON bs.region = bps.region
		AND bs.battle_id = bps.battle_id
		WHERE bs.total_players >= ?
		ORDER BY bs.start_time DESC
		LIMIT ? OFFSET ?
	`, region, playerName, playerCount, limit, offset).Scan(&summaries).Error

	return summaries, err
}

func (p *Postgres) GetBattlePlayerStats(region string, battleID int64) ([]BattlePlayerStats, error) {
	var stats []BattlePlayerStats
	err := p.db.Where("region = ? AND battle_id = ?", region, battleID).
		Order("kill_fame DESC").
		Find(&stats).Error

	return stats, err
}
