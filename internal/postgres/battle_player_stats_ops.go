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

func (p *Postgres) GetBattleSummariesByPlayer(region string, playerName string, limit int, offset int) ([]BattleSummary, error) {
	var summaries []BattleSummary
	err := p.db.Raw(`
		SELECT bs.*
		FROM battle_player_stats bps
		JOIN battle_summary bs
		  ON bs.region = bps.region
		 AND bs.battle_id = bps.battle_id
		WHERE bps.region = ?
		  AND bps.player_name = ?
		ORDER BY bs.start_time DESC
		LIMIT ? OFFSET ?
	`, region, playerName, limit, offset).Scan(&summaries).Error

	return summaries, err
}
