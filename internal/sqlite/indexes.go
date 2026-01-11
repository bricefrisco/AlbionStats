package sqlite

import (
	"fmt"

	"gorm.io/gorm"
)

func createIndexes(db *gorm.DB) error {
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_player_polls_error_count ON player_polls (error_count)`,
		`CREATE INDEX IF NOT EXISTS ON player_polls (next_poll_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_player_polls_region_next_poll ON player_polls (region, next_poll_at)`,
		`CREATE INDEX IF NOT EXISTS idx_player_stats_player_name_lower ON player_stats (region, lower(name))`,
	}

	for _, sql := range indexes {
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}
