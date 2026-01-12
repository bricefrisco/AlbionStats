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

type PlayerPveStats struct {
	Timestamps []int64 `json:"timestamps"`
	Total      []int64 `json:"total"`
	Royal      []int64 `json:"royal"`
	Outlands   []int64 `json:"outlands"`
	Avalon     []int64 `json:"avalon"`
	Hellgate   []int64 `json:"hellgate"`
	Corrupted  []int64 `json:"corrupted"`
	Mists      []int64 `json:"mists"`
}

type PlayerGatheringStats struct {
	Timestamps []int64 `json:"timestamps"`
	Total      []int64 `json:"total"`
	Royal      []int64 `json:"royal"`
	Outlands   []int64 `json:"outlands"`
	Avalon     []int64 `json:"avalon"`
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

func (s *Postgres) GetPlayerPveStats(region Region, playerID string) (*PlayerPveStats, error) {
	rows, err := s.db.
		Raw("SELECT ts, pve_total, pve_royal, pve_outlands, pve_avalon, pve_hellgate, pve_corrupted, pve_mists FROM player_stats_snapshots WHERE region = ? AND player_id = ? ORDER BY ts", region, playerID).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		stats      PlayerPveStats
		ts         time.Time
		total      int64
		royal      int64
		outlands   int64
		avalon     int64
		hellgate   int64
		corrupted  int64
		mists      int64
	)

	for rows.Next() {
		if err := rows.Scan(&ts, &total, &royal, &outlands, &avalon, &hellgate, &corrupted, &mists); err != nil {
			return nil, err
		}

		stats.Timestamps = append(stats.Timestamps, ts.UnixMilli())
		stats.Total = append(stats.Total, total)
		stats.Royal = append(stats.Royal, royal)
		stats.Outlands = append(stats.Outlands, outlands)
		stats.Avalon = append(stats.Avalon, avalon)
		stats.Hellgate = append(stats.Hellgate, hellgate)
		stats.Corrupted = append(stats.Corrupted, corrupted)
		stats.Mists = append(stats.Mists, mists)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *Postgres) GetPlayerGatheringStats(region Region, playerID string) (*PlayerGatheringStats, error) {
	rows, err := s.db.
		Raw("SELECT ts, gather_all_total, gather_all_royal, gather_all_outlands, gather_all_avalon FROM player_stats_snapshots WHERE region = ? AND player_id = ? ORDER BY ts", region, playerID).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		stats    PlayerGatheringStats
		ts       time.Time
		total    int64
		royal    int64
		outlands int64
		avalon   int64
	)

	for rows.Next() {
		if err := rows.Scan(&ts, &total, &royal, &outlands, &avalon); err != nil {
			return nil, err
		}

		stats.Timestamps = append(stats.Timestamps, ts.UnixMilli())
		stats.Total = append(stats.Total, total)
		stats.Royal = append(stats.Royal, royal)
		stats.Outlands = append(stats.Outlands, outlands)
		stats.Avalon = append(stats.Avalon, avalon)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &stats, nil
}
