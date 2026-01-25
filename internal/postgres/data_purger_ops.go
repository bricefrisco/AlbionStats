package postgres

import (
	"context"

	"gorm.io/gorm"
)

func (p *Postgres) PurgeOldBattleData(ctx context.Context) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DELETE FROM battle_summary
WHERE start_time < now() - interval '1 year'`).Error; err != nil {
			return err
		}
		if err := tx.Exec(`DELETE FROM battle_alliance_stats
WHERE start_time < now() - interval '1 year'`).Error; err != nil {
			return err
		}
		if err := tx.Exec(`DELETE FROM battle_guild_stats
WHERE start_time < now() - interval '1 year'`).Error; err != nil {
			return err
		}
		if err := tx.Exec(`DELETE FROM battle_player_stats
WHERE start_time < now() - interval '1 year'`).Error; err != nil {
			return err
		}
		if err := tx.Exec(`DELETE FROM battle_kills
WHERE ts < now() - interval '1 year'`).Error; err != nil {
			return err
		}
		if err := tx.Exec(`DELETE FROM battle_queue
WHERE ts < now() - interval '1 year'`).Error; err != nil {
			return err
		}
		return nil
	})
}
