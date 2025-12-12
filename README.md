# AlbionStats

![AlbionStats](https://i.imgur.com/zH6otmq.png)

[![CodeFactor](https://www.codefactor.io/repository/github/bricefrisco/albionstats/badge)](https://www.codefactor.io/repository/github/bricefrisco/albionstats) [![CI](https://github.com/bricefrisco/AlbionStats/actions/workflows/cicd.yml/badge.svg)](https://github.com/bricefrisco/AlbionStats/actions/workflows/cicd.yml) ![Raspberry Pi](https://img.shields.io/badge/Raspberry%20Pi-C51A4A?logo=raspberrypi&logoColor=white)

Albion Online player statistics tracker that continuously polls the game API to collect and store player data.

## Features

- **Real-time Player Statistics**: Continuously tracks and updates Albion Online player stats.
- **Historical Data**: Maintains player statistics history for trend analysis.

## Tech Stack

- **Go**: High-performance backend.
- **TimescaleDB**: PostgreSQL-based time-series database to house large amounts of data.
- **GORM**: ORM for database operations.
- **Raspberry Pi + 1TB SSD**: My cost-effective, cool side project!
