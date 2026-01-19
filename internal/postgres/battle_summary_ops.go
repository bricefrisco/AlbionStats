package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p *Postgres) InsertBattleSummaries(summaries []BattleSummary) error {
	if len(summaries) == 0 {
		return nil
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, summary := range summaries {
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&summary).Error; err != nil {
				return err
			}
		}
		return nil
	})
}