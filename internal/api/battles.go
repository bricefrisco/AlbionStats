package api

import (
	"albionstats/internal/postgres"
	"albionstats/internal/util"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MergedBattleResponse struct {
	Region        string                `json:"region"`
	BattleIDs     []int64               `json:"battle_ids"`
	StartTime     time.Time             `json:"start_time"`
	EndTime       time.Time             `json:"end_time"`
	TotalPlayers  int                   `json:"total_players"`
	TotalKills    int                   `json:"total_kills"`
	TotalFame     int64                 `json:"total_fame"`
	AllianceStats []*MergedAllianceStat `json:"alliance_stats"`
	GuildStats    []*MergedGuildStat    `json:"guild_stats"`
	PlayerStats   []*MergedPlayerStat   `json:"player_stats"`
}

type MergedAllianceStat struct {
	AllianceName string `json:"alliance_name"`
	PlayerCount  int32  `json:"player_count"`
	Kills        int32  `json:"kills"`
	Deaths       int32  `json:"deaths"`
	KillFame     int64  `json:"kill_fame"`
	DeathFame    int64  `json:"death_fame"`
	IP           int32  `json:"ip"`
}

type MergedGuildStat struct {
	GuildName    string  `json:"guild_name"`
	AllianceName *string `json:"alliance_name"`
	PlayerCount  int32   `json:"player_count"`
	Kills        int32   `json:"kills"`
	Deaths       int32   `json:"deaths"`
	KillFame     int64   `json:"kill_fame"`
	DeathFame    int64   `json:"death_fame"`
	IP           int32   `json:"ip"`
}

type MergedPlayerStat struct {
	PlayerName   string  `json:"player_name"`
	GuildName    *string `json:"guild_name"`
	AllianceName *string `json:"alliance_name"`
	Kills        int32   `json:"kills"`
	Deaths       int32   `json:"deaths"`
	KillFame     int64   `json:"kill_fame"`
	DeathFame    int64   `json:"death_fame"`
	IP           int32   `json:"ip"`
}

func (s *Server) battle(c *gin.Context) {
	region := c.Param("region")
	if !util.IsValidServer(region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region"})
		return
	}

	battleIDStr := c.Param("battleId")
	battleIDStrs := strings.Split(battleIDStr, ",")
	var battleIDs []int64
	for _, s := range battleIDStrs {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid battleId: " + s})
			return
		}
		battleIDs = append(battleIDs, id)
	}

	if len(battleIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No battleIds provided"})
		return
	}

	summaries, err := s.postgres.GetBattleSummariesByIDs(region, battleIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch battle summaries"})
		return
	}

	if len(summaries) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No battles found"})
		return
	}

	allianceStats, err := s.postgres.GetBattleAllianceStatsByIDs(region, battleIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch alliance stats"})
		return
	}

	guildStats, err := s.postgres.GetBattleGuildStatsByIDs(region, battleIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch guild stats"})
		return
	}

	playerStats, err := s.postgres.GetBattlePlayerStatsByIDs(region, battleIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	resp := s.mergeBattleSummaries(region, summaries)
	resp.AllianceStats = s.mergeAllianceStats(allianceStats)
	resp.GuildStats = s.mergeGuildStats(guildStats)
	resp.PlayerStats = s.mergePlayerStats(playerStats)

	c.JSON(http.StatusOK, resp)
}

func (s *Server) mergeBattleSummaries(region string, summaries []postgres.BattleSummary) MergedBattleResponse {
	resp := MergedBattleResponse{
		Region:    region,
		BattleIDs: make([]int64, 0, len(summaries)),
	}

	for i, summary := range summaries {
		resp.BattleIDs = append(resp.BattleIDs, summary.BattleID)
		resp.TotalPlayers += int(summary.TotalPlayers)
		resp.TotalKills += int(summary.TotalKills)
		resp.TotalFame += summary.TotalFame

		if i == 0 {
			resp.StartTime = summary.StartTime
			resp.EndTime = summary.EndTime
		} else {
			if summary.StartTime.Before(resp.StartTime) {
				resp.StartTime = summary.StartTime
			}
			if summary.EndTime.After(resp.EndTime) {
				resp.EndTime = summary.EndTime
			}
		}
	}

	return resp
}

func (s *Server) mergeAllianceStats(stats []postgres.BattleAllianceStats) []*MergedAllianceStat {
	mergedMap := make(map[string]*MergedAllianceStat)
	for _, stat := range stats {
		m, ok := mergedMap[stat.AllianceName]
		if !ok {
			m = &MergedAllianceStat{
				AllianceName: stat.AllianceName,
			}
			mergedMap[stat.AllianceName] = m
		}
		m.PlayerCount += stat.PlayerCount
		m.Kills += stat.Kills
		m.Deaths += stat.Deaths
		m.KillFame += stat.KillFame
		if stat.DeathFame != nil {
			m.DeathFame += *stat.DeathFame
		}
		if stat.IP != nil {
			m.IP += *stat.IP
		}
	}

	merged := make([]*MergedAllianceStat, 0, len(mergedMap))
	for _, v := range mergedMap {
		merged = append(merged, v)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].PlayerCount > merged[j].PlayerCount
	})

	return merged
}

func (s *Server) mergeGuildStats(stats []postgres.BattleGuildStats) []*MergedGuildStat {
	mergedMap := make(map[string]*MergedGuildStat)
	for _, stat := range stats {
		m, ok := mergedMap[stat.GuildName]
		if !ok {
			m = &MergedGuildStat{
				GuildName:    stat.GuildName,
				AllianceName: stat.AllianceName,
			}
			mergedMap[stat.GuildName] = m
		}
		m.PlayerCount += stat.PlayerCount
		m.Kills += stat.Kills
		m.Deaths += stat.Deaths
		m.KillFame += stat.KillFame
		if stat.DeathFame != nil {
			m.DeathFame += *stat.DeathFame
		}
		if stat.IP != nil {
			m.IP += *stat.IP
		}
	}

	merged := make([]*MergedGuildStat, 0, len(mergedMap))
	for _, v := range mergedMap {
		merged = append(merged, v)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].PlayerCount > merged[j].PlayerCount
	})

	return merged
}

func (s *Server) mergePlayerStats(stats []postgres.BattlePlayerStats) []*MergedPlayerStat {
	mergedMap := make(map[string]*MergedPlayerStat)
	for _, stat := range stats {
		m, ok := mergedMap[stat.PlayerName]
		if !ok {
			m = &MergedPlayerStat{
				PlayerName:   stat.PlayerName,
				GuildName:    stat.GuildName,
				AllianceName: stat.AllianceName,
			}
			mergedMap[stat.PlayerName] = m
		}
		m.Kills += stat.Kills
		m.Deaths += stat.Deaths
		m.KillFame += stat.KillFame
		if stat.DeathFame != nil {
			m.DeathFame += *stat.DeathFame
		}
		if stat.IP != nil {
			m.IP += *stat.IP
		}
	}

	merged := make([]*MergedPlayerStat, 0, len(mergedMap))
	for _, v := range mergedMap {
		merged = append(merged, v)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Kills > merged[j].Kills
	})

	return merged
}
