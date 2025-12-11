# Killboard poller

The killboard poller polls the [Albion Online events API](./albion_api.md#gameinfoevents) to collect all player names and IDs.

It then writes the data to the [player_state](./database_config.md#tables) table.

This is used later on by the [Player poller](./player_poller.md).
