# Database

## Player Polls

```sql
CREATE TABLE player_polls (
    region                          region_enum NOT NULL,
    player_id                       TEXT NOT NULL,
    next_poll_at                    TIMESTAMPTZ NOT NULL,
    error_count                     SMALLINT NOT NULL DEFAULT 0,
    killboard_last_activity         TIMESTAMPTZ,
    other_last_activity             TIMESTAMPTZ,
    last_activity                   TIMESTAMPTZ,
    last_poll_at                    TIMESTAMPTZ,

    PRIMARY KEY (region, player_id)
);

ALTER TABLE player_polls
ADD CONSTRAINT player_polls_killboard_check
CHECK (killboard_last_activity IS NOT NULL);

CREATE INDEX player_polls_poll_idx ON player_polls(next_poll_at DESC);
CREATE INDEX player_polls_region_poll_idx ON player_polls(region, next_poll_at DESC);
CREATE INDEX player_polls_error_idx ON player_polls(error_count DESC);
```

## Player Stats (Latest)

```sql
CREATE TABLE player_stats_latest (
    region                  region_enum NOT NULL,
    player_id               TEXT NOT NULL,
    ts                      TIMESTAMPTZ NOT NULL,

    -- Last seen timestamps
    killboard_last_activity TIMESTAMPTZ,
    other_last_activity     TIMESTAMPTZ,
    last_activity           TIMESTAMPTZ,
    
    -- Identity
    name                    TEXT NOT NULL,
    guild_id                TEXT,
    guild_name              TEXT,
    alliance_id             TEXT,
    alliance_name           TEXT,
    alliance_tag            TEXT,

    -- Fame counters
    kill_fame               BIGINT NOT NULL,
    death_fame              BIGINT NOT NULL,
    fame_ratio              DOUBLE PRECISION,

    -- PvE Fame
    pve_total               BIGINT NOT NULL,
    pve_royal               BIGINT NOT NULL,
    pve_outlands            BIGINT NOT NULL,
    pve_avalon              BIGINT NOT NULL,
    pve_hellgate            BIGINT NOT NULL,
    pve_corrupted           BIGINT NOT NULL,
    pve_mists               BIGINT NOT NULL,

    -- Gathering Fame Breakdown
    gather_fiber_total      BIGINT NOT NULL,
    gather_fiber_royal      BIGINT NOT NULL,
    gather_fiber_outlands   BIGINT NOT NULL,
    gather_fiber_avalon     BIGINT NOT NULL,

    gather_hide_total       BIGINT NOT NULL,
    gather_hide_royal       BIGINT NOT NULL,
    gather_hide_outlands    BIGINT NOT NULL,
    gather_hide_avalon      BIGINT NOT NULL,

    gather_ore_total        BIGINT NOT NULL,
    gather_ore_royal        BIGINT NOT NULL,
    gather_ore_outlands     BIGINT NOT NULL,
    gather_ore_avalon       BIGINT NOT NULL,

    gather_rock_total       BIGINT NOT NULL,
    gather_rock_royal       BIGINT NOT NULL,
    gather_rock_outlands    BIGINT NOT NULL,
    gather_rock_avalon      BIGINT NOT NULL,

    gather_wood_total       BIGINT NOT NULL,
    gather_wood_royal       BIGINT NOT NULL,
    gather_wood_outlands    BIGINT NOT NULL,
    gather_wood_avalon      BIGINT NOT NULL,

    gather_all_total        BIGINT NOT NULL,
    gather_all_royal        BIGINT NOT NULL,
    gather_all_outlands     BIGINT NOT NULL,
    gather_all_avalon       BIGINT NOT NULL,

    -- Crafting Fame Breakdown
    crafting_total          BIGINT NOT NULL,
    crafting_royal          BIGINT NOT NULL,
    crafting_outlands       BIGINT NOT NULL,
    crafting_avalon         BIGINT NOT NULL,

    -- Misc Lifetime Stats
    fishing_fame        BIGINT NOT NULL,
    farming_fame        BIGINT NOT NULL,
    crystal_league_fame BIGINT NOT NULL,

    PRIMARY KEY (region, player_id)
);

CREATE INDEX idx_psl_region_lower_name_prefix
ON player_stats_latest (
  region,
  lower(name) text_pattern_ops
);

CREATE INDEX idx_psl_region_lower_alliance_prefix
ON player_stats_latest (
    region,
    lower(alliance_name) text_pattern_ops
)
WHERE alliance_name IS NOT NULL;

CREATE INDEX idx_psl_region_lower_guild_prefix
ON player_stats_latest (
    region,
    lower(guild_name) text_pattern_ops
)
WHERE guild_name IS NOT NULL;

CREATE INDEX idx_psl_region_guild
ON player_stats_latest (region, guild_name)
WHERE guild_name IS NOT NULL;

CREATE INDEX idx_psl_region_alliance
ON player_stats_latest (region, alliance_name)
WHERE alliance_name IS NOT NULL;
```

## Player Stats (Snapshots)

```sql
CREATE TABLE player_stats_snapshots (
    region                  region_enum NOT NULL,
    player_id               TEXT NOT NULL,
    ts                      TIMESTAMPTZ NOT NULL,

    -- Last seen timestamps
    killboard_last_activity TIMESTAMPTZ,
    other_last_activity     TIMESTAMPTZ,
    last_activity           TIMESTAMPTZ,
    
    -- Identity
    name                    TEXT NOT NULL,
    guild_id                TEXT,
    guild_name              TEXT,
    alliance_id             TEXT,
    alliance_name           TEXT,
    alliance_tag            TEXT,

    -- Fame counters
    kill_fame               BIGINT NOT NULL,
    death_fame              BIGINT NOT NULL,
    fame_ratio              DOUBLE PRECISION,

    -- PvE Fame
    pve_total               BIGINT NOT NULL,
    pve_royal               BIGINT NOT NULL,
    pve_outlands            BIGINT NOT NULL,
    pve_avalon              BIGINT NOT NULL,
    pve_hellgate            BIGINT NOT NULL,
    pve_corrupted           BIGINT NOT NULL,
    pve_mists               BIGINT NOT NULL,

    -- Gathering Fame Breakdown
    gather_fiber_total      BIGINT NOT NULL,
    gather_fiber_royal      BIGINT NOT NULL,
    gather_fiber_outlands   BIGINT NOT NULL,
    gather_fiber_avalon     BIGINT NOT NULL,

    gather_hide_total       BIGINT NOT NULL,
    gather_hide_royal       BIGINT NOT NULL,
    gather_hide_outlands    BIGINT NOT NULL,
    gather_hide_avalon      BIGINT NOT NULL,

    gather_ore_total        BIGINT NOT NULL,
    gather_ore_royal        BIGINT NOT NULL,
    gather_ore_outlands     BIGINT NOT NULL,
    gather_ore_avalon       BIGINT NOT NULL,

    gather_rock_total       BIGINT NOT NULL,
    gather_rock_royal       BIGINT NOT NULL,
    gather_rock_outlands    BIGINT NOT NULL,
    gather_rock_avalon      BIGINT NOT NULL,

    gather_wood_total       BIGINT NOT NULL,
    gather_wood_royal       BIGINT NOT NULL,
    gather_wood_outlands    BIGINT NOT NULL,
    gather_wood_avalon      BIGINT NOT NULL,

    gather_all_total        BIGINT NOT NULL,
    gather_all_royal        BIGINT NOT NULL,
    gather_all_outlands     BIGINT NOT NULL,
    gather_all_avalon       BIGINT NOT NULL,

    -- Crafting Fame Breakdown
    crafting_total          BIGINT NOT NULL,
    crafting_royal          BIGINT NOT NULL,
    crafting_outlands       BIGINT NOT NULL,
    crafting_avalon         BIGINT NOT NULL,

    -- Misc Lifetime Stats
    fishing_fame        BIGINT NOT NULL,
    farming_fame        BIGINT NOT NULL,
    crystal_league_fame BIGINT NOT NULL,

    PRIMARY KEY (region, player_id, ts)
);

SELECT create_hypertable('player_stats_snapshots', 'ts');

ALTER TABLE player_stats_snapshots
SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'player_id',
    timescaledb.compress_orderby = 'ts DESC'
);

SELECT add_compression_policy('player_stats_snapshots', INTERVAL '1 day');
```

## Metrics

```sql
CREATE TABLE metrics (
    metric TEXT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    value BIGINT NOT NULL,
    PRIMARY KEY (metric, ts)
);

SELECT create_hypertable('metrics', 'ts');

ALTER TABLE metrics
SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'metric',
    timescaledb.compress_orderby = 'ts DESC'
);

SELECT add_compression_policy('metrics', INTERVAL '1 day');
```
