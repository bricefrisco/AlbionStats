# Albion API

## /api/gameinfo/battles

GET <https://gameinfo.albiononline.com/api/gameinfo/battles?offset=0&limit=51&sort=recent>

```json
[
    {
        "id": 1333027997,
        "startTime": "2026-01-18T21:04:33.306304400Z",
        "endTime": "2026-01-18T21:04:36.780155900Z",
        "timeout": "2026-01-18T21:07:36.780155900Z",
        "totalFame": 61596,
        "totalKills": 3,
        "clusterName": null,
        "players": {
            "LxMe-hQmSLytBgxbskhJUQ": {
                "name": "Vorthen",
                "kills": 0,
                "deaths": 1,
                "killFame": 9612,
                "guildName": "Orden Del Loto Verde",
                "guildId": "O1UMUz4QTSGOmTFXl8FPuQ",
                "allianceName": "PSS",
                "allianceId": "5u_XMWdqQY2UmaMFHPh37Q",
                "id": "LxMe-hQmSLytBgxbskhJUQ"
            },
            "hAaab9HyR8-azw5Z_JpJJw": {
                "name": "Lisa888",
                "kills": 2,
                "deaths": 0,
                "killFame": 42372,
                "guildName": "Black Lazy Wolf",
                "guildId": "gBaTlHubQ1K7Gr1dbWSxPg",
                "allianceName": "",
                "allianceId": "",
                "id": "hAaab9HyR8-azw5Z_JpJJw"
            },
            "LHJoALxMSfezZylaGoAiFA": {
                "name": "Batallita13",
                "kills": 1,
                "deaths": 1,
                "killFame": 9612,
                "guildName": "Orden Del Loto Verde",
                "guildId": "O1UMUz4QTSGOmTFXl8FPuQ",
                "allianceName": "PSS",
                "allianceId": "5u_XMWdqQY2UmaMFHPh37Q",
                "id": "LHJoALxMSfezZylaGoAiFA"
            },
            "T72jSamJTGiJDd_OmewVbA": {
                "name": "Andrjuha",
                "kills": 0,
                "deaths": 1,
                "killFame": 0,
                "guildName": "Black Lazy Wolf",
                "guildId": "gBaTlHubQ1K7Gr1dbWSxPg",
                "allianceName": "",
                "allianceId": "",
                "id": "T72jSamJTGiJDd_OmewVbA"
            }
        },
        "guilds": {
            "gBaTlHubQ1K7Gr1dbWSxPg": {
                "name": "Black Lazy Wolf",
                "kills": 2,
                "deaths": 1,
                "killFame": 42372,
                "alliance": "",
                "allianceId": "",
                "id": "gBaTlHubQ1K7Gr1dbWSxPg"
            },
            "O1UMUz4QTSGOmTFXl8FPuQ": {
                "name": "Orden Del Loto Verde",
                "kills": 1,
                "deaths": 2,
                "killFame": 19224,
                "alliance": "PSS",
                "allianceId": "5u_XMWdqQY2UmaMFHPh37Q",
                "id": "O1UMUz4QTSGOmTFXl8FPuQ"
            }
        },
        "alliances": {
            "5u_XMWdqQY2UmaMFHPh37Q": {
                "name": "PSS",
                "kills": 1,
                "deaths": 2,
                "killFame": 19224,
                "id": "5u_XMWdqQY2UmaMFHPh37Q"
            }
        },
        "battle_TIMEOUT": 180
    },
    ...
]
```
