package sqlite

import "time"

type PlayerPoll struct {
	Region                string     `gorm:"column:region;primaryKey"`
	PlayerID              string     `gorm:"column:player_id;primaryKey"`
	NextPollAt            time.Time  `gorm:"column:next_poll_at;not null"`
	ErrorCount            int        `gorm:"column:error_count;default:0"`
	LastEncountered       *time.Time `gorm:"column:last_encountered"`
	KillboardLastActivity *time.Time `gorm:"column:killboard_last_activity"`
	OtherLastActivity     *time.Time `gorm:"column:other_last_activity"`
	LastPollAt            *time.Time `gorm:"column:last_poll_at"`
}

func (PlayerPoll) TableName() string {
	return "player_polls"
}

type PlayerStats struct {
	Region   string    `gorm:"column:region;primaryKey"`
	PlayerID string    `gorm:"column:player_id;primaryKey"`
	TS       time.Time `gorm:"column:ts;not null"`

	// Last seen timestamps
	LastEncountered       *time.Time `gorm:"column:last_encountered"`
	KillboardLastActivity *time.Time `gorm:"column:killboard_last_activity"`
	OtherLastActivity     *time.Time `gorm:"column:other_last_activity"`

	// Identity
	Name         string  `gorm:"column:name;not null"`
	GuildID      *string `gorm:"column:guild_id"`
	GuildName    *string `gorm:"column:guild_name"`
	AllianceID   *string `gorm:"column:alliance_id"`
	AllianceName *string `gorm:"column:alliance_name"`
	AllianceTag  *string `gorm:"column:alliance_tag"`

	// Fame counters
	KillFame  int64   `gorm:"column:kill_fame;not null"`
	DeathFame int64   `gorm:"column:death_fame;not null"`
	FameRatio float64 `gorm:"column:fame_ratio"`

	// PvE Fame
	PveTotal     int64 `gorm:"column:pve_total;not null"`
	PveRoyal     int64 `gorm:"column:pve_royal;not null"`
	PveOutlands  int64 `gorm:"column:pve_outlands;not null"`
	PveAvalon    int64 `gorm:"column:pve_avalon;not null"`
	PveHellgate  int64 `gorm:"column:pve_hellgate;not null"`
	PveCorrupted int64 `gorm:"column:pve_corrupted;not null"`
	PveMists     int64 `gorm:"column:pve_mists;not null"`

	// Gathering Fame Breakdown
	GatherFiberTotal    int64 `gorm:"column:gather_fiber_total;not null"`
	GatherFiberRoyal    int64 `gorm:"column:gather_fiber_royal;not null"`
	GatherFiberOutlands int64 `gorm:"column:gather_fiber_outlands;not null"`
	GatherFiberAvalon   int64 `gorm:"column:gather_fiber_avalon;not null"`

	GatherHideTotal    int64 `gorm:"column:gather_hide_total;not null"`
	GatherHideRoyal    int64 `gorm:"column:gather_hide_royal;not null"`
	GatherHideOutlands int64 `gorm:"column:gather_hide_outlands;not null"`
	GatherHideAvalon   int64 `gorm:"column:gather_hide_avalon;not null"`

	GatherOreTotal    int64 `gorm:"column:gather_ore_total;not null"`
	GatherOreRoyal    int64 `gorm:"column:gather_ore_royal;not null"`
	GatherOreOutlands int64 `gorm:"column:gather_ore_outlands;not null"`
	GatherOreAvalon   int64 `gorm:"column:gather_ore_avalon;not null"`

	GatherRockTotal    int64 `gorm:"column:gather_rock_total;not null"`
	GatherRockRoyal    int64 `gorm:"column:gather_rock_royal;not null"`
	GatherRockOutlands int64 `gorm:"column:gather_rock_outlands;not null"`
	GatherRockAvalon   int64 `gorm:"column:gather_rock_avalon;not null"`

	GatherWoodTotal    int64 `gorm:"column:gather_wood_total;not null"`
	GatherWoodRoyal    int64 `gorm:"column:gather_wood_royal;not null"`
	GatherWoodOutlands int64 `gorm:"column:gather_wood_outlands;not null"`
	GatherWoodAvalon   int64 `gorm:"column:gather_wood_avalon;not null"`

	GatherAllTotal    int64 `gorm:"column:gather_all_total;not null"`
	GatherAllRoyal    int64 `gorm:"column:gather_all_royal;not null"`
	GatherAllOutlands int64 `gorm:"column:gather_all_outlands;not null"`
	GatherAllAvalon   int64 `gorm:"column:gather_all_avalon;not null"`

	// Crafting Fame Breakdown
	CraftingTotal    int64 `gorm:"column:crafting_total;not null"`
	CraftingRoyal    int64 `gorm:"column:crafting_royal;not null"`
	CraftingOutlands int64 `gorm:"column:crafting_outlands;not null"`
	CraftingAvalon   int64 `gorm:"column:crafting_avalon;not null"`

	// Misc Lifetime Stats
	FishingFame       int64 `gorm:"column:fishing_fame;not null"`
	FarmingFame       int64 `gorm:"column:farming_fame;not null"`
	CrystalLeagueFame int64 `gorm:"column:crystal_league_fame;not null"`
}

func (PlayerStats) TableName() string {
	return "player_stats"
}
