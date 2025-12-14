package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BulkUpsertStates performs a bulk upsert of player states with conflict resolution.
// Updates existing records when the new API timestamp is more recent than the stored one.
func BulkUpsertStates(ctx context.Context, db *gorm.DB, states []PlayerState) error {
	if len(states) == 0 {
		return nil
	}

	always := map[string]interface{}{
		"name":          gorm.Expr("excluded.name"),
		"guild_id":      gorm.Expr("excluded.guild_id"),
		"guild_name":    gorm.Expr("excluded.guild_name"),
		"alliance_id":   gorm.Expr("excluded.alliance_id"),
		"alliance_name": gorm.Expr("excluded.alliance_name"),
		"alliance_tag":  gorm.Expr("excluded.alliance_tag"),
		"kill_fame":     gorm.Expr("excluded.kill_fame"),
		"death_fame":    gorm.Expr("excluded.death_fame"),
		"fame_ratio":    gorm.Expr("excluded.fame_ratio"),
		"last_polled":   gorm.Expr("excluded.last_polled"),
		"last_seen":     gorm.Expr("excluded.last_seen"),
		"next_poll_at":  gorm.Expr("excluded.next_poll_at"),
		"priority":      gorm.Expr("excluded.priority"),
		"error_count":   gorm.Expr("excluded.error_count"),
		"last_error":    gorm.Expr("excluded.last_error"),
	}

	statCols := []string{
		"pve_total", "pve_royal", "pve_outlands", "pve_avalon",
		"pve_hellgate", "pve_corrupted", "pve_mists",
		"gather_fiber_total", "gather_fiber_royal", "gather_fiber_outlands", "gather_fiber_avalon",
		"gather_hide_total", "gather_hide_royal", "gather_hide_outlands", "gather_hide_avalon",
		"gather_ore_total", "gather_ore_royal", "gather_ore_outlands", "gather_ore_avalon",
		"gather_rock_total", "gather_rock_royal", "gather_rock_outlands", "gather_rock_avalon",
		"gather_wood_total", "gather_wood_royal", "gather_wood_outlands", "gather_wood_avalon",
		"gather_all_total", "gather_all_royal", "gather_all_outlands", "gather_all_avalon",
		"crafting_total", "crafting_royal", "crafting_outlands", "crafting_avalon",
		"fishing_fame", "farming_fame", "crystal_league_fame",
	}

	assignments := make(map[string]interface{}, len(always)+len(statCols))
	for k, v := range always {
		assignments[k] = v
	}
	for _, col := range statCols {
		assignments[col] = gorm.Expr(
			fmt.Sprintf(
				"CASE WHEN excluded.last_seen > player_state.last_seen THEN excluded.%s ELSE player_state.%s END",
				col, col,
			),
		)
	}

	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.Assignments(assignments),
	}).Create(&states).Error
}

// BulkUpsertFailures performs a bulk upsert of player failure states.
// Only updates error_count and next_poll_at fields.
func BulkUpsertFailures(ctx context.Context, db *gorm.DB, states []PlayerState) error {
	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"error_count",
			"next_poll_at",
		}),
	}).Create(&states).Error
}

// BulkUpsertSkips performs a bulk upsert of player skip states.
// Only updates next_poll_at and priority fields.
func BulkUpsertSkips(ctx context.Context, db *gorm.DB, states []PlayerState) error {
	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"next_poll_at",
			"priority",
		}),
	}).Create(&states).Error
}

// UpsertPlayers performs a bulk upsert of players discovered from killboard events.
// Only updates identity fields (name, guild, alliance) if the player hasn't been polled recently.
func UpsertPlayers(ctx context.Context, db *gorm.DB, players map[string]PlayerState) error {
	batch := make([]PlayerState, 0, len(players))
	for _, pl := range players {
		player := pl // copy to avoid reference issues
		batch = append(batch, player)
	}

	condition := "COALESCE(player_state.last_polled, '-infinity') <= NOW() - interval '6 hours'"

	assignments := map[string]interface{}{
		"name":          gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.name ELSE player_state.name END", condition)),
		"guild_id":      gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.guild_id ELSE player_state.guild_id END", condition)),
		"guild_name":    gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.guild_name ELSE player_state.guild_name END", condition)),
		"alliance_id":   gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.alliance_id ELSE player_state.alliance_id END", condition)),
		"alliance_name": gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.alliance_name ELSE player_state.alliance_name END", condition)),
		"alliance_tag":  gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.alliance_tag ELSE player_state.alliance_tag END", condition)),
		"last_seen":     gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.last_seen ELSE player_state.last_seen END", condition)),
		"priority":      gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN 100 ELSE player_state.priority END", condition)),
		"next_poll_at":  gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.next_poll_at ELSE player_state.next_poll_at END", condition)),
	}

	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.Assignments(assignments),
	}).Create(&batch).Error
}
