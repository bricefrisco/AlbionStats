# Database

## Battle Summary

```sql
CREATE TABLE battle_summary (
  region           region_enum,
  battle_id        BIGINT,
  start_time       TIMESTAMPTZ NOT NULL,
  end_time         TIMESTAMPTZ,

  total_players    INT NOT NULL,
  total_kills      INT NOT NULL,
  total_fame       BIGINT NOT NULL,

  alliance_names     TEXT[],
  guild_names        TEXT[],
  player_names       TEXT[],

  PRIMARY KEY(region, battle_id)
);

CREATE INDEX idx_bs_region_start_desc_inc_players
ON battle_summary (region, start_time DESC)
INCLUDE (total_players); -- Used for join with battle_player_stats

CREATE INDEX idx_bs_region_start_time_battle
ON battle_summary (region, start_time DESC, battle_id);
```

## Battle Alliance Stats

```sql
CREATE TABLE battle_alliance_stats (
  region         region_enum,
  battle_id      BIGINT,
  alliance_name  TEXT,
  start_time     TIMESTAMPTZ,
  player_count   INT,
  kills          INT,
  deaths         INT,
  kill_fame      BIGINT,

  -- Appended later, may be initially null:
  death_fame     BIGINT,
  ip             INT,

  PRIMARY KEY (region, battle_id, alliance_name)
);

CREATE INDEX idx_bas_alliance_players_battle
ON battle_alliance_stats (region, alliance_name, player_count DESC, battle_id);

CREATE INDEX idx_bgs_region_time_alliance
ON battle_alliance_stats (region, start_time, alliance_name);
```

## Battle Guild Stats

```sql
CREATE TABLE battle_guild_stats (
  region         region_enum,
  battle_id      BIGINT,
  guild_name     TEXT,
  alliance_name  TEXT,
  start_time     TIMESTAMPTZ,
  player_count   INT,
  kills          INT,
  deaths         INT,
  kill_fame      BIGINT,

  -- Appended later, may be initially null:
  death_fame     BIGINT,
  ip             INT,

  PRIMARY KEY (region, battle_id, guild_name)
);

CREATE INDEX idx_bas_guild_players_battle
ON battle_guild_stats (region, guild_name, player_count DESC, battle_id);

CREATE INDEX idx_bgs_region_time_guild
ON battle_guild_stats (region, start_time, guild_name);

CREATE INDEX idx_bgs_alliance_time_guild
ON battle_guild_stats (region, alliance_name, start_time, guild_name);
```

## Battle Player Stats

```sql
CREATE TABLE battle_player_stats (
  region         region_enum,
  battle_id      BIGINT,
  player_name    TEXT,
  start_time     TIMESTAMPTZ,
  guild_name     TEXT,
  alliance_name  TEXT,
  kills          INT,
  deaths         INT,
  kill_fame      BIGINT,

  -- Appended later, may be initially null:
  death_fame     BIGINT,
  ip             INT,
  weapon         TEXT,
  damage         BIGINT,
  heal           BIGINT,

  PRIMARY KEY (region, battle_id, player_name)
);

CREATE INDEX idx_bas_players_battle
ON battle_player_stats (region, player_name, battle_id);

CREATE INDEX idx_bgs_region_time_player
ON battle_player_stats (region, start_time, player_name);
```

## Battle queue

```sql
CREATE TABLE battle_queue (
  region         region_enum,
  battle_id      BIGINT,
  ts             TIMESTAMPTZ,
  error_count    SMALLINT NOT NULL DEFAULT 0,
  processed      BOOLEAN NOT NULL DEFAULT FALSE,

  PRIMARY KEY (region, battle_id)
);

CREATE INDEX idx_battle_queue_unprocessed_ts
ON battle_queue (ts)
WHERE processed = FALSE;
```

## Battle Kills

```sql
CREATE TABLE battle_kills (
  region         region_enum,
  battle_id      BIGINT,
  ts             TIMESTAMPTZ,
  killer_name    TEXT,
  killer_ip      INT,
  killer_weapon  TEXT,
  victim_name    TEXT,
  victim_ip      INT,
  victim_weapon  TEXT,
  fame           BIGINT
);

CREATE INDEX idx_battle_kills_region_battle_ts_desc
ON battle_kills (region, battle_id, ts DESC);
```