package postgres

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p *Postgres) InsertBattleKills(kills []BattleKills) error {
	if len(kills) == 0 {
		return nil
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, kill := range kills {
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&kill).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Postgres) GetBattleKillsByIDs(ctx context.Context, region string, battleIDs []int64) ([]BattleKills, error) {
	var kills []BattleKills
	err := p.db.WithContext(ctx).
		Where("region = ? AND battle_id IN ?", region, battleIDs).
		Order("ts DESC").
		Find(&kills).Error
	return kills, err
}
