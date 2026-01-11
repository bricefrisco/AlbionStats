package postgres

func (s *Postgres) InsertMetrics(metrics []Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	return s.db.Create(&metrics).Error
}
