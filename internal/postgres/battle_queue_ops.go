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

func (p *Postgres) GetBattleQueuesByRegion(region Region, limit int) ([]BattleQueue, error) {
	var queues []BattleQueue
	err := p.db.
		Where("region = ?", region).
		Order("ts ASC").
		Limit(limit).
		Find(&queues).Error
	return queues, err
}

func (p *Postgres) MarkBattleQueueProcessed(region Region, battleID int64) error {
	return p.db.Model(&BattleQueue{}).Where("region = ? AND battle_id = ?", region, battleID).Update("processed", true).Error
}