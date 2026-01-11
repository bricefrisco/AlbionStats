package postgres

func (s *Postgres) InsertPlayerStatsSnapshots(stats []PlayerStatsSnapshot) error {
	if len(stats) == 0 {
		return nil
	}

	return s.db.Create(&stats).Error
}
