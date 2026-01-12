package postgres

import (
	"database/sql"
	"time"
)

func (s *Postgres) InsertPlayerStatsSnapshots(stats []PlayerStatsSnapshot) error {
	if len(stats) == 0 {
		return nil
	}

	return s.db.Create(&stats).Error
}

type PlayerPvpStats struct {
	Timestamps []int64     `json:"timestamps"`
	KillFame   []int64     `json:"kill_fame"`
	DeathFame  []int64     `json:"death_fame"`
	FameRatio  []*float64  `json:"fame_ratio"`
}

func (s *Postgres) GetPlayerPvpStats(region Region, playerID string) (*PlayerPvpStats, error) {
	rows, err := s.db.
		Raw("SELECT ts, kill_fame, death_fame, fame_ratio FROM player_stats_snapshots WHERE region = ? AND player_id = ? ORDER BY ts", region, playerID).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		stats     PlayerPvpStats
		ts        time.Time
		killFame  int64
		deathFame int64
		fameRatio sql.NullFloat64
	)

	for rows.Next() {
		if err := rows.Scan(&ts, &killFame, &deathFame, &fameRatio); err != nil {
			return nil, err
		}

		stats.Timestamps = append(stats.Timestamps, ts.UnixMilli())
		stats.KillFame = append(stats.KillFame, killFame)
		stats.DeathFame = append(stats.DeathFame, deathFame)

		if fameRatio.Valid {
			value := fameRatio.Float64
			stats.FameRatio = append(stats.FameRatio, &value)
		} else {
			stats.FameRatio = append(stats.FameRatio, nil)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &stats, nil
}
