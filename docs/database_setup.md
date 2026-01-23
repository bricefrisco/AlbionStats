# Database configuration

Database was set up using [instructions from timescale db](https://www.tigerdata.com/docs/self-hosted/latest/install/installation-linux)

```cmd
sudo apt install gnupg postgresql-common apt-transport-https lsb-release wget
sudo /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh
echo "deb https://packagecloud.io/timescale/timescaledb/ubuntu/ $(lsb_release -c -s) main" | sudo tee /etc/apt/sources.list.d/timescaledb.list
wget --quiet -O - https://packagecloud.io/timescale/timescaledb/gpgkey | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/timescaledb.gpg
sudo apt update
sudo apt install timescaledb-2-postgresql-18 postgresql-client-18
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
sudo nano /etc/postgresql/18/main/postgresql.conf
listen_addresses = '*'
```

```cmd
sudo nano /etc/postgresql/18/main/pg_hba.conf
host    all             all             127.0.0.1/32            scram-sha-256
host    all             all             X.X.X.X/32              scram-sha-256
sudo systemctl restart postgresql
```

```sql
CREATE EXTENSION IF NOT EXISTS timescaledb;
```

```cmd
\dx
```

## Region Enum

```sql
CREATE TYPE region_enum AS ENUM ( 'americas', 'europe', 'asia');
```
