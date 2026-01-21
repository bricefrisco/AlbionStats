# Database Backups

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
FILE=${DB_NAME}_${DATE}.dump.gz

mkdir -p "$BACKUP_DIR"

log "Starting backup"
pg_dump -Fc "$DB_NAME" | gzip > "$BACKUP_DIR/$FILE"
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

```cmd
crontab -e
```

```
CRON_TZ=UTC
0 11 * * 2 sudo -u postgres /usr/local/bin/pg_backup.sh
```


