package sqlite

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *SQLite) FetchPlayersToPoll(ctx context.Context, region string, batchSize int) ([]PlayerPoll, error) {
	var players []PlayerPoll
	now := time.Now().UTC()
	if err := s.db.WithContext(ctx).
		Where("region = ? AND next_poll_at <= ?", region, now).
		Order("next_poll_at ASC").
		Limit(batchSize).
		Find(&players).Error; err != nil {
		return nil, fmt.Errorf("query players: %w", err)
	}
	return players, nil
}

func (s *SQLite) UpsertPlayerPolls(polls map[string]PlayerPoll) error {
	if len(polls) == 0 {
		return nil
	}

	batch := make([]PlayerPoll, 0, len(polls))
	for _, poll := range polls {
		batch = append(batch, poll)
	}

	assignmentsOnConflict := map[string]interface{}{
		"next_poll_at": gorm.Expr(
			"MIN(" +
				"COALESCE(datetime(player_polls.last_poll_at, '+6 hours'), player_polls.next_poll_at)," +
				"player_polls.next_poll_at" +
				")",
		),
		"killboard_last_activity": gorm.Expr("excluded.killboard_last_activity"),
	}

	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.Assignments(assignmentsOnConflict),
	}).Create(&batch).Error
}

func (s *SQLite) UpdatePlayerPolls(ctx context.Context, db *gorm.DB, polls []PlayerPoll) error {
	if len(polls) == 0 {
		return nil
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, p := range polls {
			if err := tx.Model(&PlayerPoll{}).
				Where("region = ? AND player_id = ?", p.Region, p.PlayerID).
				Updates(map[string]any{
					"next_poll_at":            p.NextPollAt,
					"error_count":             p.ErrorCount,
					"other_last_activity":     p.OtherLastActivity,
					"last_poll_at":            p.LastPollAt,
					"last_encountered":        p.LastEncountered,
					"killboard_last_activity": p.KillboardLastActivity,
				}).Error; err != nil {
				return fmt.Errorf("update player poll: %w", err)
			}
		}
		return nil
	})
}

func (s *SQLite) UpsertPlayerStats(ctx context.Context, db *gorm.DB, stats []PlayerStats) error {
	if len(stats) == 0 {
		return nil
	}

	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		UpdateAll: true,
	}).Create(&stats).Error
}

// func ApplyPlayerPollerDatabaseChanges(ctx context.Context, db *gorm.DB, deletePolls []PlayerPoll, updatePolls []PlayerPoll, stats []PlayerStats, failures []PlayerPoll) error {
// 	if len(deletePolls) > 0 {
// 		regionBuckets := make(map[string][]string)
// 		for _, d := range deletePolls {
// 			regionBuckets[d.Region] = append(regionBuckets[d.Region], d.PlayerID)
// 		}
// 		for region, ids := range regionBuckets {
// 			if err := db.WithContext(ctx).Delete(&PlayerPoll{}, "region = ? AND player_id IN ?", region, ids).Error; err != nil {
// 				return fmt.Errorf("delete polls: %w", err)
// 			}
// 		}
// 	}

// 	// Upsert successful polls
// 	if len(updatePolls) > 0 {
// 		if err := BulkUpsertPlayerPolls(ctx, db, updatePolls); err != nil {
// 			return fmt.Errorf("upsert polls: %w", err)
// 		}
// 	}

// 	// Upsert stats latest
// 	if len(stats) > 0 {
// 		if err := BulkUpsertPlayerStats(ctx, db, statsLatest); err != nil {
// 			return fmt.Errorf("upsert stats latest: %w", err)
// 		}
// 	}

// 	// Insert snapshots
// 	if len(snapshots) > 0 {
// 		if err := BulkInsertPlayerStatsSnapshots(ctx, db, snapshots); err != nil {
// 			return fmt.Errorf("insert snapshots: %w", err)
// 		}
// 	}

// 	// Upsert failures
// 	if len(failurePolls) > 0 {
// 		if err := BulkUpsertPlayerPolls(ctx, db, failurePolls); err != nil {
// 			return fmt.Errorf("upsert failure polls: %w", err)
// 		}
// 	}

// 	return nil
// }
