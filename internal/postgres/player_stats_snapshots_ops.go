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

type PlayerPvpSeries struct {
	KillFame  []int64    `json:"kill_fame"`
	DeathFame []int64    `json:"death_fame"`
	FameRatio []*float64 `json:"fame_ratio"`
}

type PlayerPveSeries struct {
	Total     []int64 `json:"total"`
	Royal     []int64 `json:"royal"`
	Outlands  []int64 `json:"outlands"`
	Avalon    []int64 `json:"avalon"`
	Hellgate  []int64 `json:"hellgate"`
	Corrupted []int64 `json:"corrupted"`
	Mists     []int64 `json:"mists"`
}

type PlayerGatheringSeries struct {
	Total    []int64 `json:"total"`
	Royal    []int64 `json:"royal"`
	Outlands []int64 `json:"outlands"`
	Avalon   []int64 `json:"avalon"`
}

type PlayerCraftingSeries struct {
	Total []int64 `json:"total"`
}

type PlayerStatsSeries struct {
	Timestamps []int64               `json:"timestamps"`
	Pvp        PlayerPvpSeries       `json:"pvp"`
	Pve        PlayerPveSeries       `json:"pve"`
	Gathering  PlayerGatheringSeries `json:"gathering"`
	Crafting   PlayerCraftingSeries  `json:"crafting"`
}

func (s *Postgres) GetPlayerStatsSeries(region Region, playerID string) (*PlayerStatsSeries, error) {
	rows, err := s.db.
		Raw("SELECT ts, kill_fame, death_fame, fame_ratio, pve_total, pve_royal, pve_outlands, pve_avalon, pve_hellgate, pve_corrupted, pve_mists, gather_all_total, gather_all_royal, gather_all_outlands, gather_all_avalon, crafting_total FROM player_stats_snapshots WHERE region = ? AND player_id = ? ORDER BY ts", region, playerID).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		series     PlayerStatsSeries
		ts         time.Time
		killFame   int64
		deathFame  int64
		fameRatio  sql.NullFloat64
		pveTotal   int64
		pveRoyal   int64
		pveOut     int64
		pveAvalon  int64
		pveHell    int64
		pveCorr    int64
		pveMists   int64
		gathTotal  int64
		gathRoyal  int64
		gathOut    int64
		gathAval   int64
		craftTotal int64
	)

	for rows.Next() {
		if err := rows.Scan(
			&ts,
			&killFame,
			&deathFame,
			&fameRatio,
			&pveTotal,
			&pveRoyal,
			&pveOut,
			&pveAvalon,
			&pveHell,
			&pveCorr,
			&pveMists,
			&gathTotal,
			&gathRoyal,
			&gathOut,
			&gathAval,
			&craftTotal,
		); err != nil {
			return nil, err
		}

		series.Timestamps = append(series.Timestamps, ts.UnixMilli())
		series.Pvp.KillFame = append(series.Pvp.KillFame, killFame)
		series.Pvp.DeathFame = append(series.Pvp.DeathFame, deathFame)
		if fameRatio.Valid {
			value := fameRatio.Float64
			series.Pvp.FameRatio = append(series.Pvp.FameRatio, &value)
		} else {
			series.Pvp.FameRatio = append(series.Pvp.FameRatio, nil)
		}

		series.Pve.Total = append(series.Pve.Total, pveTotal)
		series.Pve.Royal = append(series.Pve.Royal, pveRoyal)
		series.Pve.Outlands = append(series.Pve.Outlands, pveOut)
		series.Pve.Avalon = append(series.Pve.Avalon, pveAvalon)
		series.Pve.Hellgate = append(series.Pve.Hellgate, pveHell)
		series.Pve.Corrupted = append(series.Pve.Corrupted, pveCorr)
		series.Pve.Mists = append(series.Pve.Mists, pveMists)

		series.Gathering.Total = append(series.Gathering.Total, gathTotal)
		series.Gathering.Royal = append(series.Gathering.Royal, gathRoyal)
		series.Gathering.Outlands = append(series.Gathering.Outlands, gathOut)
		series.Gathering.Avalon = append(series.Gathering.Avalon, gathAval)

		series.Crafting.Total = append(series.Crafting.Total, craftTotal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &series, nil
}
