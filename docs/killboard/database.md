# Database

## Battle Summary

```sql
CREATE TABLE battle_summary (
  battle_id        BIGINT PRIMARY KEY,
  start_time       TIMESTAMPTZ NOT NULL,
  end_time         TIMESTAMPTZ,
  cluster_name     TEXT,

  total_players    INT NOT NULL,
  total_kills      INT NOT NULL,
  total_deaths     INT NOT NULL,
  total_fame       BIGINT NOT NULL,

  alliance_names     TEXT[],
  guild_names        TEXT[],
  player_names       TEXT[]
);

CREATE INDEX ON battle_summary (start_time DESC);
CREATE INDEX ON battle_summary USING GIN ((lower(guild_names)));
CREATE INDEX ON battle_summary USING GIN ((lower(alliance_names)));
CREATE INDEX ON battle_summary USING GIN ((lower(player_names)));
```

## Battle Alliance Stats

```sql
CREATE TABLE battle_alliance_stats (
  battle_id      BIGINT,
  alliance_name  TEXT,

  player_count   INT,
  kills          INT,
  deaths         INT,
  avg_ip         INT,
  total_fame     BIGINT,

  PRIMARY KEY (battle_id, alliance_name)
);
```

## Battle Guild Stats

```sql
CREATE TABLE battle_guild_stats (
  battle_id      BIGINT,
  guild_name     TEXT,
  alliance_name  TEXT,
  player_count   INT,
  kills          INT,
  deaths         INT,
  avg_ip         INT,
  total_fame     BIGINT,

  PRIMARY KEY (battle_id, guild_name)
);
```

## Battle Player Stats

```sql
CREATE TABLE battle_player_stats (
  battle_id      BIGINT,
  player_name    TEXT,
  guild_name     TEXT,
  alliance_name  TEXT,

  weapon         TEXT,
  ip             INT,

  damage_done    BIGINT,
  healing_done   BIGINT,
  kills          INT,
  deaths         INT,
  total_fame     BIGINT,

  PRIMARY KEY (battle_id, player_name)
);

SELECT create_hypertable(
  'battle_player_stats',
  by_range('battle_id')
);

ALTER TABLE battle_player_stats SET (
  timescaledb.compress,
  timescaledb.compress_segmentby = 'battle_id'
);
```

## Battle Kills

```sql
CREATE TABLE battle_kills (
  battle_id      BIGINT,
  ts             TIMESTAMPTZ,

  killer_id      TEXT,
  killer_name    TEXT,
  killer_ip      INT,

  victim_id      TEXT,
  victim_name    TEXT,
  victim_ip      INT,

  fame           BIGINT
);

SELECT create_hypertable(
  'battle_kills',
  by_range('ts')
);

ALTER TABLE battle_kills SET (
  timescaledb.compress,
  timescaledb.compress_segmentby = 'battle_id'
);
```
