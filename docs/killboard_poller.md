# Killboard poller

The killboard poller polls the [Albion Online events API](./albion_api.md#gameinfoevents) to collect all player names and IDs.

It then writes the data to the [player_polls](./database_config.md#tables) table using:

```sql
INSERT INTO player_polls (region, player_id, first_seen, next_poll_at, killboard_last_activity)
VALUES (?, ?, TRUE, NOW(), ?)
ON CONFLICT (region, player_id)
DO UPDATE SET
    killboard_last_activity = EXCLUDED.killboard_last_activity,
    next_poll_at = LEAST(player_polls.last_poll_at + INTERVAL '6 hours', player_polls.next_poll_at)
```

This is used later on by the [Player poller](./player_poller.md).
