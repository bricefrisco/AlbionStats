package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TopPlayerStats struct {
	PlayerName     string `gorm:"column:player_name"`
	TotalKillFame  int64  `gorm:"column:total_kill_fame"`
	TotalDeathFame int64  `gorm:"column:total_death_fame"`
	TotalKills     int64  `gorm:"column:total_kills"`
	TotalDeaths    int64  `gorm:"column:total_deaths"`
}

type AlliancePlayerStats struct {
	PlayerName string    `gorm:"column:player_name"`
	LastBattle time.Time `gorm:"column:last_battle"`
	NumBattles int64     `gorm:"column:num_battles"`
	Kills      int64     `gorm:"column:kills"`
	Deaths     int64     `gorm:"column:deaths"`
	KillFame   int64     `gorm:"column:kill_fame"`
	DeathFame  int64     `gorm:"column:death_fame"`
}

func (p *Postgres) InsertBattlePlayerStats(stats []BattlePlayerStats) error {
	if len(stats) == 0 {
		return nil
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&stat).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Postgres) UpdateBattlePlayerStats(stats []BattlePlayerStats) error {
	if len(stats) == 0 {
		return nil
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			updates := make(map[string]interface{})
			updates["death_fame"] = stat.DeathFame
			updates["ip"] = stat.IP
			updates["weapon"] = stat.Weapon
			updates["damage"] = stat.Damage
			updates["heal"] = stat.Heal

			if err := tx.Model(&BattlePlayerStats{}).
				Where("region = ? AND battle_id = ? AND player_name = ?", stat.Region, stat.BattleID, stat.PlayerName).
				Updates(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Postgres) GetBattleSummariesByPlayer(region string, playerName string, playerCount int, limit int, offset int) ([]BattleSummary, error) {
	var summaries []BattleSummary
	err := p.db.Raw(`
		SELECT bs.*
		FROM (
		    SELECT region, battle_id
		    FROM battle_player_stats
		    WHERE region = ?
		    AND player_name = ?
		) bps
		JOIN battle_summary bs
		ON bs.region = bps.region
		AND bs.battle_id = bps.battle_id
		WHERE bs.total_players >= ?
		ORDER BY bs.start_time DESC
		LIMIT ? OFFSET ?
	`, region, playerName, playerCount, limit, offset).Scan(&summaries).Error

	return summaries, err
}

func (p *Postgres) GetBattlePlayerStatsByIDs(ctx context.Context, region string, battleIDs []int64) ([]BattlePlayerStats, error) {
	var stats []BattlePlayerStats
	err := p.db.WithContext(ctx).Where("region = ? AND battle_id IN ?", region, battleIDs).
		Order("kills DESC").
		Find(&stats).Error

	return stats, err
}

func (p *Postgres) GetAlliancePlayerStats(region string, allianceName string) ([]AlliancePlayerStats, error) {
	var stats []AlliancePlayerStats
	err := p.db.Raw(`
		SELECT
			bps.player_name,
			MAX(bps.start_time) AS last_battle,
			COUNT(DISTINCT bps.battle_id) AS num_battles,
			SUM(bps.kills) AS kills,
			SUM(bps.deaths) AS deaths,
			SUM(bps.kill_fame) AS kill_fame,
			SUM(bps.death_fame) AS death_fame
		FROM battle_player_stats bps
		WHERE bps.region = ?
			AND bps.alliance_name = ?
			AND bps.start_time >= NOW() - INTERVAL '30 days'
		GROUP BY bps.player_name
		ORDER BY kill_fame DESC
	`, region, allianceName).Scan(&stats).Error

	return stats, err
}

func (p *Postgres) GetTopPlayers(region string, limit int, offset int) ([]TopPlayerStats, error) {
	var stats []TopPlayerStats
	err := p.db.Raw(`
		SELECT
			player_name,
			SUM(kill_fame) AS total_kill_fame,
			SUM(COALESCE(death_fame, 0)) AS total_death_fame,
			SUM(kills) AS total_kills,
			SUM(deaths) AS total_deaths
		FROM battle_player_stats
		WHERE region = ?
			AND start_time >= NOW() - INTERVAL '30 days'
		GROUP BY player_name
		ORDER BY total_kill_fame DESC
		LIMIT ? OFFSET ?
	`, region, limit, offset).Scan(&stats).Error

	return stats, err
}
