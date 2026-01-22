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

func (p *Postgres) GetBattleSummariesByGuild(region string, guildName string, playerCount int, limit int, offset int) ([]BattleSummary, error) {
	var summaries []BattleSummary
	err := p.db.Raw(`
		SELECT bs.*
		FROM (
		    SELECT region, battle_id
		    FROM battle_guild_stats
		    WHERE region = ?
		    AND guild_name = ?
		    AND player_count >= ?
		) bgs
		JOIN battle_summary bs
		ON bs.region = bgs.region
		AND bs.battle_id = bgs.battle_id
		ORDER BY bs.start_time DESC
		LIMIT ? OFFSET ?
	`, region, guildName, playerCount, limit, offset).Scan(&summaries).Error

	return summaries, err
}
