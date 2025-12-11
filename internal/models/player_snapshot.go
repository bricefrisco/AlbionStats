package models

import "time"

// PlayerStatsSnapshot maps to player_stats_snapshots table.
type PlayerStatsSnapshot struct {
	Region       string     `gorm:"column:region;primaryKey;type:region_enum"`
	PlayerID     string     `gorm:"column:player_id;primaryKey"`
	TS           time.Time  `gorm:"column:ts;primaryKey"`
	APITimestamp *time.Time `gorm:"column:api_timestamp"`

	// Identity
	Name         string  `gorm:"column:name"`
	GuildID      *string `gorm:"column:guild_id"`
	GuildName    *string `gorm:"column:guild_name"`
	AllianceID   *string `gorm:"column:alliance_id"`
	AllianceName *string `gorm:"column:alliance_name"`
	AllianceTag  *string `gorm:"column:alliance_tag"`

	// Fame counters
	KillFame  *int64   `gorm:"column:kill_fame"`
	DeathFame *int64   `gorm:"column:death_fame"`
	FameRatio *float64 `gorm:"column:fame_ratio"`

	// PvE Fame
	PveTotal     *int64 `gorm:"column:pve_total"`
	PveRoyal     *int64 `gorm:"column:pve_royal"`
	PveOutlands  *int64 `gorm:"column:pve_outlands"`
	PveAvalon    *int64 `gorm:"column:pve_avalon"`
	PveHellgate  *int64 `gorm:"column:pve_hellgate"`
	PveCorrupted *int64 `gorm:"column:pve_corrupted"`
	PveMists     *int64 `gorm:"column:pve_mists"`

	// Gathering
	GatherFiberTotal    *int64 `gorm:"column:gather_fiber_total"`
	GatherFiberRoyal    *int64 `gorm:"column:gather_fiber_royal"`
	GatherFiberOutlands *int64 `gorm:"column:gather_fiber_outlands"`
	GatherFiberAvalon   *int64 `gorm:"column:gather_fiber_avalon"`

	GatherHideTotal    *int64 `gorm:"column:gather_hide_total"`
	GatherHideRoyal    *int64 `gorm:"column:gather_hide_royal"`
	GatherHideOutlands *int64 `gorm:"column:gather_hide_outlands"`
	GatherHideAvalon   *int64 `gorm:"column:gather_hide_avalon"`

	GatherOreTotal    *int64 `gorm:"column:gather_ore_total"`
	GatherOreRoyal    *int64 `gorm:"column:gather_ore_royal"`
	GatherOreOutlands *int64 `gorm:"column:gather_ore_outlands"`
	GatherOreAvalon   *int64 `gorm:"column:gather_ore_avalon"`

	GatherRockTotal    *int64 `gorm:"column:gather_rock_total"`
	GatherRockRoyal    *int64 `gorm:"column:gather_rock_royal"`
	GatherRockOutlands *int64 `gorm:"column:gather_rock_outlands"`
	GatherRockAvalon   *int64 `gorm:"column:gather_rock_avalon"`

	GatherWoodTotal    *int64 `gorm:"column:gather_wood_total"`
	GatherWoodRoyal    *int64 `gorm:"column:gather_wood_royal"`
	GatherWoodOutlands *int64 `gorm:"column:gather_wood_outlands"`
	GatherWoodAvalon   *int64 `gorm:"column:gather_wood_avalon"`

	GatherAllTotal    *int64 `gorm:"column:gather_all_total"`
	GatherAllRoyal    *int64 `gorm:"column:gather_all_royal"`
	GatherAllOutlands *int64 `gorm:"column:gather_all_outlands"`
	GatherAllAvalon   *int64 `gorm:"column:gather_all_avalon"`

	// Crafting
	CraftingTotal    *int64 `gorm:"column:crafting_total"`
	CraftingRoyal    *int64 `gorm:"column:crafting_royal"`
	CraftingOutlands *int64 `gorm:"column:crafting_outlands"`
	CraftingAvalon   *int64 `gorm:"column:crafting_avalon"`

	// Misc
	FishingFame       *int64 `gorm:"column:fishing_fame"`
	FarmingFame       *int64 `gorm:"column:farming_fame"`
	CrystalLeagueFame *int64 `gorm:"column:crystal_league_fame"`
}

func (PlayerStatsSnapshot) TableName() string {
	return "player_stats_snapshots"
}
