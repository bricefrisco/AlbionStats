# AlbionStats

### https://albionstats.com

<img src="https://i.imgur.com/dqYYN9p.png" alt="AlbionStats" width="100">

[![CI](https://github.com/bricefrisco/AlbionStats/actions/workflows/cicd.yml/badge.svg)](https://github.com/bricefrisco/AlbionStats/actions/workflows/cicd.yml)

Albion Online player statistics tracker that continuously polls the game API to collect and store data.

## What it does

- **Ingests Albion Online game data**: Polls players, killboard events, battles, and battleboards for Americas, Europe, and Asia.
- **Auto-tracks players from killboards**: When a player appears in killboard events, they get queued for tracking.
- **Stores rich player history**: Keeps latest stats plus time-series snapshots for PvE, PvP, gathering, and crafting.
- **Builds battle insights**: Aggregates battle summaries, alliance/guild/player stats, and kill logs.
- **Collects platform metrics**: Tracks total players, total snapshots/data points, and daily active users by region.
- **Maintains data health**: Periodically purges old battle data.

## Tech Stack

- **Go**: High-performance backend.
- **SvelteKit**: Web frontend.
- **Gin**: HTTP API server.
- **PostgreSQL + TimescaleDB**: relational and time-series database to house data.
- **GORM**: ORM for database operations.
