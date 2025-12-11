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
(Modify line) listen_addresses = '192.168.1.192'
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

## Tables

```sql
CREATE TYPE region_enum AS ENUM ( 'americas', 'europe', 'asia');

CREATE TABLE player_state (
    region                  region_enum NOT NULL,
    player_id               TEXT NOT NULL,

    -- Identity
    name                    TEXT,
    guild_id                TEXT,
    guild_name              TEXT,
    alliance_id             TEXT,
    alliance_name           TEXT,
    alliance_tag            TEXT,

    -- Activity tracking
    last_seen               TIMESTAMPTZ,
    last_polled             TIMESTAMPTZ,

    -- Scheduling
    priority                INTEGER NOT NULL DEFAULT 300,
    next_poll_at            TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Error handling
    error_count             INTEGER NOT NULL DEFAULT 0,
    last_error              TEXT,

    -- Fame counters
    kill_fame               BIGINT,
    death_fame              BIGINT,
    fame_ratio              DOUBLE PRECISION,

    -- PvE Fame
    pve_total               BIGINT,
    pve_royal               BIGINT,
    pve_outlands            BIGINT,
    pve_avalon              BIGINT,
    pve_hellgate            BIGINT,
    pve_corrupted           BIGINT,
    pve_mists               BIGINT,

    -- Gathering Fame Breakdown
    gather_fiber_total      BIGINT,
    gather_fiber_royal      BIGINT,
    gather_fiber_outlands   BIGINT,
    gather_fiber_avalon     BIGINT,

    gather_hide_total       BIGINT,
    gather_hide_royal       BIGINT,
    gather_hide_outlands    BIGINT,
    gather_hide_avalon      BIGINT,

    gather_ore_total        BIGINT,
    gather_ore_royal        BIGINT,
    gather_ore_outlands     BIGINT,
    gather_ore_avalon       BIGINT,

    gather_rock_total       BIGINT,
    gather_rock_royal       BIGINT,
    gather_rock_outlands    BIGINT,
    gather_rock_avalon      BIGINT,

    gather_wood_total       BIGINT,
    gather_wood_royal       BIGINT,
    gather_wood_outlands    BIGINT,
    gather_wood_avalon      BIGINT,

    gather_all_total        BIGINT,
    gather_all_royal        BIGINT,
    gather_all_outlands     BIGINT,
    gather_all_avalon       BIGINT,

    -- Crafting Fame Breakdown
    crafting_total          BIGINT,
    crafting_royal          BIGINT,
    crafting_outlands       BIGINT,
    crafting_avalon         BIGINT,

    -- Misc Lifetime Stats
    fishing_fame        BIGINT,
    farming_fame        BIGINT,
    crystal_league_fame BIGINT,

    PRIMARY KEY (region, player_id)
);

CREATE TABLE player_stats_snapshots (
    region                  region_enum NOT NULL,
    player_id               TEXT NOT NULL,
    ts                      TIMESTAMPTZ NOT NULL,
    api_timestamp           TIMESTAMPTZ,

    -- Identity
    name                    TEXT,
    guild_id                TEXT,
    guild_name              TEXT,
    alliance_id             TEXT,
    alliance_name           TEXT,
    alliance_tag            TEXT,

    -- Fame counters
    kill_fame               BIGINT,
    death_fame              BIGINT,
    fame_ratio              DOUBLE PRECISION,

    -- PvE Fame
    pve_total               BIGINT,
    pve_royal               BIGINT,
    pve_outlands            BIGINT,
    pve_avalon              BIGINT,
    pve_hellgate            BIGINT,
    pve_corrupted           BIGINT,
    pve_mists               BIGINT,

    -- Gathering Fame Breakdown
    gather_fiber_total      BIGINT,
    gather_fiber_royal      BIGINT,
    gather_fiber_outlands   BIGINT,
    gather_fiber_avalon     BIGINT,

    gather_hide_total       BIGINT,
    gather_hide_royal       BIGINT,
    gather_hide_outlands    BIGINT,
    gather_hide_avalon      BIGINT,

    gather_ore_total        BIGINT,
    gather_ore_royal        BIGINT,
    gather_ore_outlands     BIGINT,
    gather_ore_avalon       BIGINT,

    gather_rock_total       BIGINT,
    gather_rock_royal       BIGINT,
    gather_rock_outlands    BIGINT,
    gather_rock_avalon      BIGINT,

    gather_wood_total       BIGINT,
    gather_wood_royal       BIGINT,
    gather_wood_outlands    BIGINT,
    gather_wood_avalon      BIGINT,

    gather_all_total        BIGINT,
    gather_all_royal        BIGINT,
    gather_all_outlands     BIGINT,
    gather_all_avalon       BIGINT,

    -- Crafting Fame Breakdown
    crafting_total          BIGINT,
    crafting_royal          BIGINT,
    crafting_outlands       BIGINT,
    crafting_avalon         BIGINT,

    -- Misc Lifetime Stats
    fishing_fame        BIGINT,
    farming_fame        BIGINT,
    crystal_league_fame BIGINT,

    PRIMARY KEY (region, player_id, ts)
);

SELECT create_hypertable('player_stats_snapshots', 'ts');
```
