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

CREATE INDEX ON battle_summary (region, start_time DESC);
```

## Battle Alliance Stats

```sql
CREATE TABLE battle_alliance_stats (
  region         region_enum,
  battle_id      BIGINT,
  alliance_name  TEXT,
  player_count   INT,
  kills          INT,
  deaths         INT,
  kill_fame      BIGINT,

  -- Appended later, may be initially null:
  death_fame     BIGINT,
  ip             INT,

  PRIMARY KEY (region, battle_id, alliance_name)
);
```

## Battle Guild Stats

```sql
CREATE TABLE battle_guild_stats (
  region         region_enum,
  battle_id      BIGINT,
  guild_name     TEXT,
  alliance_name  TEXT,
  player_count   INT,
  kills          INT,
  deaths         INT,
  kill_fame      BIGINT,

  -- Appended later, may be initially null:
  death_fame     BIGINT,
  ip             INT,

  PRIMARY KEY (region, battle_id, guild_name)
);
```

## Battle Player Stats

```sql
CREATE TABLE battle_player_stats (
  region         region_enum,
  battle_id      BIGINT,
  player_name    TEXT,
  guild_name     TEXT,
  alliance_name  TEXT,
  kills          INT,
  deaths         INT,
  kill_fame      BIGINT,

  -- Appended later, may be initially null:
  death_fame     BIGINT,
  ip             INT,

  PRIMARY KEY (region, battle_id, player_name)
);
```

## Battle queue

```sql
CREATE TABLE battle_queue (
  region         region_enum,
  battle_id      BIGINT,
  ts             TIMESTAMPTZ,
  error_count    SMALLINT NOT NULL DEFAULT 0,

  PRIMARY KEY (region, battle_id)
);

CREATE INDEX ON battle_queue (ts ASC);
```

## Battle Kills

```sql
CREATE TABLE battle_kills (
  region         region_enum,
  battle_id      BIGINT,
  ts             TIMESTAMPTZ,
  killer_name    TEXT,
  killer_ip      INT,
  victim_name    TEXT,
  victim_ip      INT,
  fame           BIGINT
);

CREATE INDEX ON battle_kills (region, battle_id, ts DESC);
```
