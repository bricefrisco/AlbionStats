package postgres

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Postgres) FetchPlayersToPoll(region string, batchSize int) ([]PlayerPoll, error) {
	var players []PlayerPoll
	now := time.Now().UTC()
	if err := s.db.Where("region = ? AND next_poll_at <= ?", region, now).
		Order("next_poll_at ASC").
		Limit(batchSize).
		Find(&players).Error; err != nil {
		return nil, fmt.Errorf("query players: %w", err)
	}
	return players, nil
}

func (s *Postgres) UpsertPlayerPolls(polls map[string]PlayerPoll) error {
	if len(polls) == 0 {
		return nil
	}

	batch := make([]PlayerPoll, 0, len(polls))
	for _, poll := range polls {
		batch = append(batch, poll)
	}

	assignmentsOnConflict := map[string]interface{}{
		"next_poll_at": gorm.Expr(
			"LEAST(" +
				"COALESCE(player_polls.last_poll_at + INTERVAL '12 hours', player_polls.next_poll_at)," +
				"player_polls.next_poll_at" +
				")",
		),
		"killboard_last_activity": gorm.Expr("excluded.killboard_last_activity"),
		"last_activity":           gorm.Expr("excluded.last_activity"),
	}

	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.Assignments(assignmentsOnConflict),
	}).Create(&batch).Error
}

func (s *Postgres) UpdatePlayerPolls(polls []PlayerPoll) error {
	if len(polls) == 0 {
		return nil
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, p := range polls {
			if err := tx.Model(&PlayerPoll{}).
				Where("region = ? AND player_id = ?", p.Region, p.PlayerID).
				Updates(map[string]any{
					"next_poll_at":            p.NextPollAt,
					"error_count":             p.ErrorCount,
					"other_last_activity":     p.OtherLastActivity,
					"last_poll_at":            p.LastPollAt,
					"last_activity":           p.LastActivity,
					"killboard_last_activity": p.KillboardLastActivity,
				}).Error; err != nil {
				return fmt.Errorf("update player poll: %w", err)
			}
		}
		return nil
	})
}

func (s *Postgres) DeletePlayerPolls(polls []PlayerPoll) error {
	if len(polls) == 0 {
		return nil
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, p := range polls {
			if err := tx.Model(&PlayerPoll{}).
				Where("region = ? AND player_id = ?", p.Region, p.PlayerID).
				Delete(&PlayerPoll{}).Error; err != nil {
				return fmt.Errorf("delete player poll: %w", err)
			}
		}
		return nil
	})
}

func (s *Postgres) GetPlayersReadyToPollCount() (int64, error) {
	var count int64
	now := time.Now().UTC()
	if err := s.db.Model(&PlayerPoll{}).
		Where("next_poll_at <= ?", now).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count players ready to poll: %w", err)
	}
	return count, nil
}

func (s *Postgres) GetPlayersWithErrorsCount() (int64, error) {
	var count int64
	if err := s.db.Model(&PlayerPoll{}).
		Where("error_count >= ?", 1).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count players with errors: %w", err)
	}
	return count, nil
}
