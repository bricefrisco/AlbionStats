package postgres

import "gorm.io/gorm/clause"

func (s *Postgres) UpsertPlayerStatsLatest(stats []PlayerStatsLatest) error {
	if len(stats) == 0 {
		return nil
	}

	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		UpdateAll: true,
	}).Create(&stats).Error
}
