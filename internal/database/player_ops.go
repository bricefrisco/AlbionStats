package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpsertKillboardPlayerPolls(ctx context.Context, db *gorm.DB, polls map[string]PlayerPoll) error {
	if len(polls) == 0 {
		return nil
	}

	batch := make([]PlayerPoll, 0, len(polls))
	for _, poll := range polls {
		playerPoll := poll // copy to avoid reference issues
		batch = append(batch, playerPoll)
	}

	assignments := map[string]interface{}{
		"next_poll_at":            gorm.Expr("LEAST(player_polls.last_poll_at + INTERVAL '6 hours', player_polls.next_poll_at)"),
		"killboard_last_activity": gorm.Expr("excluded.killboard_last_activity"),
	}

	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.Assignments(assignments),
	}).Create(&batch).Error
}

func UpsertLastEncounteredPlayerPolls(ctx context.Context, db *gorm.DB, poll PlayerPoll) error {
	assignments := map[string]interface{}{
		"next_poll_at":     gorm.Expr("LEAST(player_polls.last_poll_at + INTERVAL '6 hours', player_polls.next_poll_at)"),
		"last_encountered": gorm.Expr("excluded.last_encountered"),
	}

	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.Assignments(assignments),
	}).Create(&poll).Error
}

// BulkUpsertPlayerStatsLatest performs a bulk upsert of player stats latest records.
// Always updates the most recent stats for each player.
func BulkUpsertPlayerStatsLatest(ctx context.Context, db *gorm.DB, stats []PlayerStatsLatest) error {
	if len(stats) == 0 {
		return nil
	}

	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		UpdateAll: true,
	}).Create(&stats).Error
}

// BulkInsertPlayerStatsSnapshots performs a bulk insert of player stats snapshots.
// Uses TimescaleDB hypertable for time-series data.
func BulkInsertPlayerStatsSnapshots(ctx context.Context, db *gorm.DB, snapshots []PlayerStatsSnapshot) error {
	if len(snapshots) == 0 {
		return nil
	}

	return db.WithContext(ctx).Create(&snapshots).Error
}

// BulkUpsertPlayerPollsKillboard performs a bulk upsert of player polls.
// Updates all relevant fields for both successful polls and failures.
func BulkUpsertPlayerPolls(ctx context.Context, db *gorm.DB, polls []PlayerPoll) error {
	if len(polls) == 0 {
		return nil
	}

	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"next_poll_at", "error_count", "other_last_activity", "last_poll_at",
			"last_encountered", "killboard_last_activity",
		}),
	}).Create(&polls).Error
}

// ApplyPlayerPollerDatabaseChanges applies all database changes from player polling in bulk.
func ApplyPlayerPollerDatabaseChanges(ctx context.Context, db *gorm.DB, deletePolls []PlayerPoll, updatePolls []PlayerPoll, statsLatest []PlayerStatsLatest, snapshots []PlayerStatsSnapshot, failurePolls []PlayerPoll) error {
	// Apply deletes
	if len(deletePolls) > 0 {
		regionBuckets := make(map[string][]string)
		for _, d := range deletePolls {
			regionBuckets[d.Region] = append(regionBuckets[d.Region], d.PlayerID)
		}
		for region, ids := range regionBuckets {
			if err := db.WithContext(ctx).Delete(&PlayerPoll{}, "region = ? AND player_id IN ?", region, ids).Error; err != nil {
				return fmt.Errorf("delete polls: %w", err)
			}
		}
	}

	// Upsert successful polls
	if len(updatePolls) > 0 {
		if err := BulkUpsertPlayerPolls(ctx, db, updatePolls); err != nil {
			return fmt.Errorf("upsert polls: %w", err)
		}
	}

	// Upsert stats latest
	if len(statsLatest) > 0 {
		if err := BulkUpsertPlayerStatsLatest(ctx, db, statsLatest); err != nil {
			return fmt.Errorf("upsert stats latest: %w", err)
		}
	}

	// Insert snapshots
	if len(snapshots) > 0 {
		if err := BulkInsertPlayerStatsSnapshots(ctx, db, snapshots); err != nil {
			return fmt.Errorf("insert snapshots: %w", err)
		}
	}

	// Upsert failures
	if len(failurePolls) > 0 {
		if err := BulkUpsertPlayerPolls(ctx, db, failurePolls); err != nil {
			return fmt.Errorf("upsert failure polls: %w", err)
		}
	}

	return nil
}
