# Database configuration

Database was set up using [instructions from timescale db](https://www.tigerdata.com/docs/self-hosted/latest/install/installation-linux)

```cmd
sudo apt install gnupg postgresql-common apt-transport-https lsb-release wget
sudo /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh
echo "deb https://packagecloud.io/timescale/timescaledb/debian/ $(lsb_release -c -s) main" | sudo tee /etc/apt/sources.list.d/timescaledb.list
wget --quiet -O - https://packagecloud.io/timescale/timescaledb/gpgkey | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/timescaledb.gpg
sudo apt update
sudo apt install timescaledb-2-postgresql-17 postgresql-client-17
sudo timescaledb-tune
```

(Accepted all default options)

```cmd
sudo systemctl restart postgresql
sudo -u postgres psql
\password postgres
```

Allowing local connections:

```cmd
sudo nano /etc/postgresql/17/main/postgresql.conf
(Modify line) listen_addresses = 'localhost,192.168.1.192'
```

```cmd
sudo nano /etc/postgresql/17/main/pg_hba.conf
host    all    all    192.168.1.0/24    md5
(Add line)
sudo systemctl restart postgresql
```

```sql
CREATE EXTENSION IF NOT EXISTS timescaledb;
```

```cmd
\dx
```

## Data Access Patterns

## Tables

```sql
CREATE TYPE region_enum AS ENUM ( 'americas', 'europe', 'asia');

-----------------------
-- player_polls
-----------------------
CREATE TABLE player_polls (
    region                          region_enum NOT NULL,
    player_id                       TEXT NOT NULL,
    next_poll_at                    TIMESTAMPTZ NOT NULL,
    error_count                     SMALLINT NOT NULL DEFAULT 0,
    last_encountered                TIMESTAMPTZ,
    killboard_last_activity         TIMESTAMPTZ,
    other_last_activity             TIMESTAMPTZ,
    last_poll_at                    TIMESTAMPTZ,

    PRIMARY KEY (region, player_id),

    CHECK (
        last_encountered IS NOT NULL OR     -- Added through the API
        killboard_last_activity IS NOT NULL -- Added through killboard crawler
    )
);

CREATE INDEX player_polls_poll_idx ON player_polls(region, next_poll_at DESC);

-----------------------
-- player_stats_latest
-----------------------
CREATE TABLE player_stats_latest (
    region                  region_enum NOT NULL,
    player_id               TEXT NOT NULL,
    ts                      TIMESTAMPTZ NOT NULL,

    -- Last seen timestamps
    last_encountered        TIMESTAMPTZ,
    killboard_last_activity TIMESTAMPTZ,
    other_last_activity     TIMESTAMPTZ,
    
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

CREATE INDEX idx_player_name_lower
ON player_stats_latest (region, lower(name));

-----------------------
-- player_stats_snapshots
-----------------------
CREATE TABLE player_stats_snapshots (
    region                  region_enum NOT NULL,
    player_id               TEXT NOT NULL,
    ts                      TIMESTAMPTZ NOT NULL,

    -- Last seen timestamps
    last_encountered        TIMESTAMPTZ,
    killboard_last_activity TIMESTAMPTZ,
    other_last_activity     TIMESTAMPTZ,
    
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

-----------------------
-- metrics
-----------------------
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
