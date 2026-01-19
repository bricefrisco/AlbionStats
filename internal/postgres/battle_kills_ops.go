package postgres

import (
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
