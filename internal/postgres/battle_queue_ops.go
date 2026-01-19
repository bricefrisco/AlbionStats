package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p *Postgres) InsertBattleQueues(queues []BattleQueue) error {
	if len(queues) == 0 {
		return nil
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, queue := range queues {
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&queue).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
