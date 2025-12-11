package models

import "time"

// PlayerState maps to the player_state table.
type PlayerState struct {
	Region       string     `gorm:"column:region;primaryKey;type:region_enum"`
	PlayerID     string     `gorm:"column:player_id;primaryKey"`
	Name         string     `gorm:"column:name"`
	GuildID      *string    `gorm:"column:guild_id"`
	GuildName    *string    `gorm:"column:guild_name"`
	AllianceID   *string    `gorm:"column:alliance_id"`
	AllianceName *string    `gorm:"column:alliance_name"`
	AllianceTag  *string    `gorm:"column:alliance_tag"`
	LastSeen     *time.Time `gorm:"column:last_seen"`
	LastPolled   *time.Time `gorm:"column:last_polled"`
	Priority     int        `gorm:"column:priority;default:300"`
	NextPollAt   time.Time  `gorm:"column:next_poll_at;default:now()"`
	ErrorCount   int        `gorm:"column:error_count;default:0"`
	LastError    *string    `gorm:"column:last_error"`
}

func (PlayerState) TableName() string {
	return "player_state"
}
