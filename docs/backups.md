# Backups

## Unzip
Needed to run rclone.
```cmd
apt update
apt install -y unzip
```

## Rclone
Used to sync backups to Backblaze B2.
```cmd
curl https://rclone.org/install.sh | bash
rclone version
sudo -u postgres rclone config
```

## Backup directory
```cmd
mkdir -p /var/backups/postgres
chown postgres:postgres /var/backups/postgres
chmod 700 /var/backups/postgres
```

## Backup logs
```cmd
touch /var/log/pg_backup.log
chown postgres:postgres /var/log/pg_backup.log
chmod 600 /var/log/pg_backup.log
```

## Script

`nano /usr/local/bin/pg_backup.sh`

```bash
#!/bin/bash
set -euo pipefail

LOG_FILE=/var/log/pg_backup.log
exec >>"$LOG_FILE" 2>&1

log() { echo "$(date '+%Y-%m-%d %H:%M:%S') | $1"; }

DB_NAME=postgres
BACKUP_DIR=/var/backups/postgres
REMOTE=b2:pg-backups
DATE=$(date +%Y-%m-%d)
FILE=${DB_NAME}_${DATE}.bak

mkdir -p "$BACKUP_DIR"

log "Starting backup"
pg_dump -Fc "$DB_NAME" -f "$BACKUP_DIR/$FILE"
log "Upload to B2"
rclone copy "$BACKUP_DIR/$FILE" "$REMOTE" --checksum
log "Delete older remote backups"
rclone delete "$REMOTE" --exclude "$FILE"
log "Remove local file"
rm -f "$BACKUP_DIR/$FILE"
log "Done"
```

```cmd
chmod +x /usr/local/bin/pg_backup.sh
```

## Running script
```cmd
sudo -u postgres /usr/local/bin/pg_backup.sh
tail -f /var/log/pg_backup.log
```

## Cron

```cmd
crontab -e
```

```
CRON_TZ=UTC
0 11 * * 2 sudo -u postgres /usr/local/bin/pg_backup.sh
```

# Restore

```cmd
sudo su postgres
dropdb <database>
createdb <database>
```

```cmd
psql -d <database>
\c <database>
SELECT timescaledb_pre_restore();
\q
```

```cmd
pg_restore -Fc -d <database> <file>.bak
psql -d <database>
\c <database>
SELECT timescaledb_post_restore();
\q
```

## Decompression

For some reason, compressed data was not visible upon restore.

Recompressing resolved the issue:
```cmd
SELECT decompress_chunk(format('%I.%I', chunk_schema, chunk_name)::regclass)
FROM timescaledb_information.chunks
WHERE hypertable_name = 'player_stats_snapshots'
  AND is_compressed = true;
 
ALTER TABLE player_stats_snapshots
SET (timescaledb.compress = false);
  
SELECT compress_chunk(format('%I.%I', chunk_schema, chunk_name)::regclass)
FROM timescaledb_information.chunks
WHERE hypertable_name = 'player_stats_snapshots'
  AND is_compressed = false;

ALTER TABLE player_stats_snapshots
SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'player_id',
    timescaledb.compress_orderby = 'ts DESC'
);
```

The `metrics` hypertable was unaffected.
