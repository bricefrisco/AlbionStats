package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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