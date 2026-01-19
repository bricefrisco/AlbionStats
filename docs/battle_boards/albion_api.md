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

## /api/gameinfo/events/battle/{id}

GET <https://gameinfo.albiononline.com/api/gameinfo/events/battle/1333329647?offset=0&limit=51>

```json
[
  {
    "groupMemberCount": 3,
    "numberOfParticipants": 2,
    "EventId": 1333329647,
    "TimeStamp": "2026-01-19T13:21:12.218997700Z",
    "Version": 4,
    "Killer": {
      "AverageItemPower": 1236.091,
      "Equipment": {
        "MainHand": {
          "Type": "T5_MAIN_FIRESTAFF_KEEPER@2",
          "Count": 1,
          "Quality": 4,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "OffHand": {
          "Type": "T6_OFF_BOOK@1",
          "Count": 1,
          "Quality": 4,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Head": {
          "Type": "T6_HEAD_CLOTH_SET2",
          "Count": 1,
          "Quality": 4,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Armor": {
          "Type": "T5_ARMOR_CLOTH_KEEPER@2",
          "Count": 1,
          "Quality": 2,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Shoes": {
          "Type": "T6_SHOES_PLATE_SET1",
          "Count": 1,
          "Quality": 3,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Bag": {
          "Type": "T5_BAG",
          "Count": 1,
          "Quality": 3,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Cape": {
          "Type": "T5_CAPEITEM_MORGANA@3",
          "Count": 1,
          "Quality": 5,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Mount": {
          "Type": "T6_MOUNT_FROSTRAM_ADC",
          "Count": 1,
          "Quality": 1,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Potion": {
          "Type": "T5_POTION_REVIVE",
          "Count": 4,
          "Quality": 0,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Food": {
          "Type": "T8_MEAL_STEW",
          "Count": 1,
          "Quality": 0,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        }
      },
      "Inventory": [],
      "Name": "flotjs",
      "Id": "6YnAPYHzSF6Kcrry5MFDgw",
      "GuildName": "Da Black Lotus",
      "GuildId": "Yiyb9T12TxeCWCz78BU7Fw",
      "AllianceName": "FR3YA",
      "AllianceId": "NB18Mi5BSTqYtkxiIE--aw",
      "AllianceTag": "",
      "Avatar": "AVATAR_AJ_CHARACTER_PROGRESSION_01",
      "AvatarRing": "AVATARRING_AJ_CHARACTER_PROGRESSION_01",
      "DeathFame": 0,
      "KillFame": 96,
      "FameRatio": 960,
      "LifetimeStatistics": {
        "PvE": {
          "Total": 0,
          "Royal": 0,
          "Outlands": 0,
          "Avalon": 0,
          "Hellgate": 0,
          "CorruptedDungeon": 0,
          "Mists": 0
        },
        "Gathering": {
          "Fiber": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "Hide": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "Ore": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "Rock": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "Wood": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "All": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          }
        },
        "Crafting": {
          "Total": 0,
          "Royal": 0,
          "Outlands": 0,
          "Avalon": 0
        },
        "CrystalLeague": 0,
        "FishingFame": 0,
        "FarmingFame": 0,
        "Timestamp": null
      }
    },
    "Victim": {
      "AverageItemPower": 1185.161,
      "Equipment": {
        "MainHand": {
          "Type": "T5_2H_BOW_KEEPER@2",
          "Count": 1,
          "Quality": 2,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "OffHand": null,
        "Head": {
          "Type": "T6_HEAD_CLOTH_SET3@1",
          "Count": 1,
          "Quality": 2,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Armor": {
          "Type": "T5_ARMOR_LEATHER_SET3@2",
          "Count": 1,
          "Quality": 1,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Shoes": {
          "Type": "T5_SHOES_PLATE_HELL@2",
          "Count": 1,
          "Quality": 2,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Bag": {
          "Type": "T6_BAG",
          "Count": 1,
          "Quality": 2,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Cape": {
          "Type": "T4_CAPEITEM_FW_LYMHURST@2",
          "Count": 1,
          "Quality": 1,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Mount": null,
        "Potion": {
          "Type": "T4_POTION_COOLDOWN",
          "Count": 9,
          "Quality": 0,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        "Food": {
          "Type": "T5_MEAL_SOUP",
          "Count": 2,
          "Quality": 0,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        }
      },
      "Inventory": [
        {
          "Type": "T4_HEAD_PLATE_SET2",
          "Count": 1,
          "Quality": 2,
          "ActiveSpells": [],
          "PassiveSpells": [],
          "LegendarySoul": null
        },
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null
      ],
      "Name": "AbraXxX",
      "Id": "cUuY2xTKRIm1IhT2HYs8dw",
      "GuildName": "Ascensa0 Bad0nica",
      "GuildId": "U0T9WdAgRn-esjetFK35sA",
      "AllianceName": "",
      "AllianceId": "",
      "AllianceTag": "",
      "Avatar": "AVATAR_05",
      "AvatarRing": "AVATARRING_ADC_APR2019",
      "DeathFame": 288,
      "KillFame": 0,
      "FameRatio": 0,
      "LifetimeStatistics": {
        "PvE": {
          "Total": 0,
          "Royal": 0,
          "Outlands": 0,
          "Avalon": 0,
          "Hellgate": 0,
          "CorruptedDungeon": 0,
          "Mists": 0
        },
        "Gathering": {
          "Fiber": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "Hide": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "Ore": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "Rock": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "Wood": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "All": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          }
        },
        "Crafting": {
          "Total": 0,
          "Royal": 0,
          "Outlands": 0,
          "Avalon": 0
        },
        "CrystalLeague": 0,
        "FishingFame": 0,
        "FarmingFame": 0,
        "Timestamp": null
      }
    },
    "TotalVictimKillFame": 288,
    "Location": null,
    "Participants": [
      {
        "AverageItemPower": 1236.091,
        "Equipment": {
          "MainHand": {
            "Type": "T5_MAIN_FIRESTAFF_KEEPER@2",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "OffHand": {
            "Type": "T6_OFF_BOOK@1",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Head": {
            "Type": "T6_HEAD_CLOTH_SET2",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Armor": {
            "Type": "T5_ARMOR_CLOTH_KEEPER@2",
            "Count": 1,
            "Quality": 2,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Shoes": {
            "Type": "T6_SHOES_PLATE_SET1",
            "Count": 1,
            "Quality": 3,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Bag": {
            "Type": "T5_BAG",
            "Count": 1,
            "Quality": 3,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Cape": {
            "Type": "T5_CAPEITEM_MORGANA@3",
            "Count": 1,
            "Quality": 5,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Mount": {
            "Type": "T6_MOUNT_FROSTRAM_ADC",
            "Count": 1,
            "Quality": 1,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Potion": {
            "Type": "T5_POTION_REVIVE",
            "Count": 4,
            "Quality": 0,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Food": {
            "Type": "T8_MEAL_STEW",
            "Count": 1,
            "Quality": 0,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          }
        },
        "Inventory": [],
        "Name": "flotjs",
        "Id": "6YnAPYHzSF6Kcrry5MFDgw",
        "GuildName": "Da Black Lotus",
        "GuildId": "Yiyb9T12TxeCWCz78BU7Fw",
        "AllianceName": "FR3YA",
        "AllianceId": "NB18Mi5BSTqYtkxiIE--aw",
        "AllianceTag": "",
        "Avatar": "AVATAR_AJ_CHARACTER_PROGRESSION_01",
        "AvatarRing": "AVATARRING_AJ_CHARACTER_PROGRESSION_01",
        "DeathFame": 0,
        "KillFame": 0,
        "FameRatio": 0,
        "LifetimeStatistics": {
          "PvE": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0,
            "Hellgate": 0,
            "CorruptedDungeon": 0,
            "Mists": 0
          },
          "Gathering": {
            "Fiber": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Hide": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Ore": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Rock": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Wood": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "All": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            }
          },
          "Crafting": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "CrystalLeague": 0,
          "FishingFame": 0,
          "FarmingFame": 0,
          "Timestamp": null
        },
        "DamageDone": 2132,
        "SupportHealingDone": 0
      },
      {
        "AverageItemPower": 1284.208,
        "Equipment": {
          "MainHand": {
            "Type": "T8_MAIN_HOLYSTAFF@3",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "OffHand": {
            "Type": "T8_OFF_HORN_KEEPER@1",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Head": {
            "Type": "T8_HEAD_CLOTH_SET1@1",
            "Count": 1,
            "Quality": 2,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Armor": {
            "Type": "T8_ARMOR_CLOTH_SET2@2",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Shoes": {
            "Type": "T6_SHOES_PLATE_SET2@3",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Bag": null,
          "Cape": {
            "Type": "T8_CAPEITEM_FW_MARTLOCK@3",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Mount": {
            "Type": "T7_MOUNT_DIREBOAR",
            "Count": 1,
            "Quality": 2,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Potion": {
            "Type": "T7_POTION_STONESKIN",
            "Count": 4,
            "Quality": 0,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "Food": {
            "Type": "T7_MEAL_OMELETTE_FISH",
            "Count": 3,
            "Quality": 0,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          }
        },
        "Inventory": [],
        "Name": "Indiera",
        "Id": "FEHBMN3fS1i3MsaAQvx8Ww",
        "GuildName": "M0RS",
        "GuildId": "3Wn-JdABSb2NplpCqiFN5Q",
        "AllianceName": "MST",
        "AllianceId": "JNzo-2DHTSqQM91lHno8lg",
        "AllianceTag": "",
        "Avatar": "AVATAR_AJ_CHARACTER_PROGRESSION_01",
        "AvatarRing": "AVATARRING_ADC_JAN2019",
        "DeathFame": 0,
        "KillFame": 0,
        "FameRatio": 0,
        "LifetimeStatistics": {
          "PvE": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0,
            "Hellgate": 0,
            "CorruptedDungeon": 0,
            "Mists": 0
          },
          "Gathering": {
            "Fiber": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Hide": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Ore": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Rock": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Wood": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "All": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            }
          },
          "Crafting": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "CrystalLeague": 0,
          "FishingFame": 0,
          "FarmingFame": 0,
          "Timestamp": null
        },
        "DamageDone": 0,
        "SupportHealingDone": 581
      }
    ],
    "GroupMembers": [
      {
        "AverageItemPower": 0,
        "Equipment": {
          "MainHand": {
            "Type": "T8_MAIN_HOLYSTAFF@3",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "OffHand": null,
          "Head": null,
          "Armor": null,
          "Shoes": null,
          "Bag": null,
          "Cape": null,
          "Mount": null,
          "Potion": null,
          "Food": null
        },
        "Inventory": [],
        "Name": "Indiera",
        "Id": "FEHBMN3fS1i3MsaAQvx8Ww",
        "GuildName": "M0RS",
        "GuildId": "3Wn-JdABSb2NplpCqiFN5Q",
        "AllianceName": "MST",
        "AllianceId": "JNzo-2DHTSqQM91lHno8lg",
        "AllianceTag": "",
        "Avatar": "AVATAR_AJ_CHARACTER_PROGRESSION_01",
        "AvatarRing": "AVATARRING_ADC_JAN2019",
        "DeathFame": 0,
        "KillFame": 96,
        "FameRatio": 960,
        "LifetimeStatistics": {
          "PvE": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0,
            "Hellgate": 0,
            "CorruptedDungeon": 0,
            "Mists": 0
          },
          "Gathering": {
            "Fiber": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Hide": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Ore": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Rock": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Wood": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "All": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            }
          },
          "Crafting": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "CrystalLeague": 0,
          "FishingFame": 0,
          "FarmingFame": 0,
          "Timestamp": null
        }
      },
      {
        "AverageItemPower": 0,
        "Equipment": {
          "MainHand": {
            "Type": "T5_MAIN_FIRESTAFF_KEEPER@2",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": null
          },
          "OffHand": null,
          "Head": null,
          "Armor": null,
          "Shoes": null,
          "Bag": null,
          "Cape": null,
          "Mount": null,
          "Potion": null,
          "Food": null
        },
        "Inventory": [],
        "Name": "flotjs",
        "Id": "6YnAPYHzSF6Kcrry5MFDgw",
        "GuildName": "Da Black Lotus",
        "GuildId": "Yiyb9T12TxeCWCz78BU7Fw",
        "AllianceName": "FR3YA",
        "AllianceId": "NB18Mi5BSTqYtkxiIE--aw",
        "AllianceTag": "",
        "Avatar": "AVATAR_AJ_CHARACTER_PROGRESSION_01",
        "AvatarRing": "AVATARRING_AJ_CHARACTER_PROGRESSION_01",
        "DeathFame": 0,
        "KillFame": 96,
        "FameRatio": 960,
        "LifetimeStatistics": {
          "PvE": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0,
            "Hellgate": 0,
            "CorruptedDungeon": 0,
            "Mists": 0
          },
          "Gathering": {
            "Fiber": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Hide": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Ore": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Rock": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Wood": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "All": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            }
          },
          "Crafting": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "CrystalLeague": 0,
          "FishingFame": 0,
          "FarmingFame": 0,
          "Timestamp": null
        }
      },
      {
        "AverageItemPower": 0,
        "Equipment": {
          "MainHand": {
            "Type": "T4_2H_BOW_KEEPER@4",
            "Count": 1,
            "Quality": 4,
            "ActiveSpells": [],
            "PassiveSpells": [],
            "LegendarySoul": {
              "id": "617fc59e-ff14-428b-adb9-769599010b06",
              "subtype": 1,
              "era": 1,
              "name": null,
              "lastEquipped": "2026-01-19T13:16:39.977688200Z",
              "attunedPlayer": "60afcfe7-541a-4042-90d1-d3808a4f3b4d",
              "attunedPlayerName": "MrPEVEZ",
              "attunement": 113484556,
              "attunementSpentSinceReset": 72501598,
              "attunementSpent": 72501598,
              "quality": 4,
              "craftedBy": "MrPEVEZ",
              "traits": [
                {
                  "roll": 0.3398074876554986,
                  "pendingRolls": [],
                  "pendingTraits": [],
                  "value": 0.0571575537,
                  "trait": "TRAIT_ABILITY_DAMAGE",
                  "minvalue": 0.00165,
                  "maxvalue": 0.165
                },
                {
                  "roll": 0.10898161207386349,
                  "pendingRolls": [],
                  "pendingTraits": [],
                  "value": 0.029472949,
                  "trait": "TRAIT_CAST_SPEED_INCREASE",
                  "minvalue": 0.0025,
                  "maxvalue": 0.25
                }
              ],
              "PvPFameGained": 960000
            }
          },
          "OffHand": null,
          "Head": null,
          "Armor": null,
          "Shoes": null,
          "Bag": null,
          "Cape": null,
          "Mount": null,
          "Potion": null,
          "Food": null
        },
        "Inventory": [],
        "Name": "KALIZEL",
        "Id": "F45aRZ1SR-q3Q-R3cwZD0g",
        "GuildName": "Mas Respeto",
        "GuildId": "nVikl6GUSSuKM7nNJgGQpA",
        "AllianceName": "HUMO",
        "AllianceId": "TCIYuZ3MQq-C_kFpZxGYbw",
        "AllianceTag": "",
        "Avatar": "AVATAR_FAMERANK_01",
        "AvatarRing": "RING1",
        "DeathFame": 0,
        "KillFame": 96,
        "FameRatio": 960,
        "LifetimeStatistics": {
          "PvE": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0,
            "Hellgate": 0,
            "CorruptedDungeon": 0,
            "Mists": 0
          },
          "Gathering": {
            "Fiber": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Hide": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Ore": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Rock": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "Wood": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            },
            "All": {
              "Total": 0,
              "Royal": 0,
              "Outlands": 0,
              "Avalon": 0
            }
          },
          "Crafting": {
            "Total": 0,
            "Royal": 0,
            "Outlands": 0,
            "Avalon": 0
          },
          "CrystalLeague": 0,
          "FishingFame": 0,
          "FarmingFame": 0,
          "Timestamp": null
        }
      }
    ],
    "GvGMatch": null,
    "BattleId": 1333329647,
    "KillArea": "OPEN_WORLD",
    "Category": null,
    "Type": "KILL"
  },
```
