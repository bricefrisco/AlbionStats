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
	Region        string
	BattleIDs     []int64
	StartTime     time.Time
	EndTime       time.Time
	TotalPlayers  int
	TotalKills    int
	TotalFame     int64
	AllianceStats []*MergedAllianceStat `json:"Alliances"`
	GuildStats    []*MergedGuildStat    `json:"Guilds"`
	PlayerStats   []*MergedPlayerStat   `json:"Players"`
}

type MergedAllianceStat struct {
	AllianceName string
	PlayerCount  int32
	Kills        int32
	Deaths       int32
	KillFame     int64
	DeathFame    int64
	IP           int32
}

type MergedGuildStat struct {
	GuildName    string
	AllianceName *string
	PlayerCount  int32
	Kills        int32
	Deaths       int32
	KillFame     int64
	DeathFame    int64
	IP           int32
}

type MergedPlayerStat struct {
	PlayerName   string
	GuildName    *string
	AllianceName *string
	Kills        int32
	Deaths       int32
	KillFame     int64
	DeathFame    int64
	IP           int32
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

	if len(battleIDs) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 20 battleIds allowed"})
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
	type ipCalc struct {
		totalWeightedIP int64
		count           int64
	}
	ipMap := make(map[string]*ipCalc)

	for _, stat := range stats {
		m, ok := mergedMap[stat.AllianceName]
		if !ok {
			m = &MergedAllianceStat{
				AllianceName: stat.AllianceName,
			}
			mergedMap[stat.AllianceName] = m
			ipMap[stat.AllianceName] = &ipCalc{}
		}
		m.PlayerCount += stat.PlayerCount
		m.Kills += stat.Kills
		m.Deaths += stat.Deaths
		m.KillFame += stat.KillFame
		if stat.DeathFame != nil {
			m.DeathFame += *stat.DeathFame
		}
		if stat.IP != nil {
			ipMap[stat.AllianceName].totalWeightedIP += int64(*stat.IP) * int64(stat.PlayerCount)
			ipMap[stat.AllianceName].count += int64(stat.PlayerCount)
		}
	}

	merged := make([]*MergedAllianceStat, 0, len(mergedMap))
	for name, v := range mergedMap {
		if ipMap[name].count > 0 {
			v.IP = int32(ipMap[name].totalWeightedIP / ipMap[name].count)
		}
		merged = append(merged, v)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].PlayerCount > merged[j].PlayerCount
	})

	return merged
}

func (s *Server) mergeGuildStats(stats []postgres.BattleGuildStats) []*MergedGuildStat {
	mergedMap := make(map[string]*MergedGuildStat)
	type ipCalc struct {
		totalWeightedIP int64
		count           int64
	}
	ipMap := make(map[string]*ipCalc)

	for _, stat := range stats {
		m, ok := mergedMap[stat.GuildName]
		if !ok {
			m = &MergedGuildStat{
				GuildName:    stat.GuildName,
				AllianceName: stat.AllianceName,
			}
			mergedMap[stat.GuildName] = m
			ipMap[stat.GuildName] = &ipCalc{}
		}
		m.PlayerCount += stat.PlayerCount
		m.Kills += stat.Kills
		m.Deaths += stat.Deaths
		m.KillFame += stat.KillFame
		if stat.DeathFame != nil {
			m.DeathFame += *stat.DeathFame
		}
		if stat.IP != nil {
			ipMap[stat.GuildName].totalWeightedIP += int64(*stat.IP) * int64(stat.PlayerCount)
			ipMap[stat.GuildName].count += int64(stat.PlayerCount)
		}
	}

	merged := make([]*MergedGuildStat, 0, len(mergedMap))
	for name, v := range mergedMap {
		if ipMap[name].count > 0 {
			v.IP = int32(ipMap[name].totalWeightedIP / ipMap[name].count)
		}
		merged = append(merged, v)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].PlayerCount > merged[j].PlayerCount
	})

	return merged
}

func (s *Server) mergePlayerStats(stats []postgres.BattlePlayerStats) []*MergedPlayerStat {
	mergedMap := make(map[string]*MergedPlayerStat)
	type ipCalc struct {
		totalIP int64
		count   int64
	}
	ipMap := make(map[string]*ipCalc)

	for _, stat := range stats {
		m, ok := mergedMap[stat.PlayerName]
		if !ok {
			m = &MergedPlayerStat{
				PlayerName:   stat.PlayerName,
				GuildName:    stat.GuildName,
				AllianceName: stat.AllianceName,
			}
			mergedMap[stat.PlayerName] = m
			ipMap[stat.PlayerName] = &ipCalc{}
		}
		m.Kills += stat.Kills
		m.Deaths += stat.Deaths
		m.KillFame += stat.KillFame
		if stat.DeathFame != nil {
			m.DeathFame += *stat.DeathFame
		}
		if stat.IP != nil {
			ipMap[stat.PlayerName].totalIP += int64(*stat.IP)
			ipMap[stat.PlayerName].count++
		}
	}

	merged := make([]*MergedPlayerStat, 0, len(mergedMap))
	for name, v := range mergedMap {
		if ipMap[name].count > 0 {
			v.IP = int32(ipMap[name].totalIP / ipMap[name].count)
		}
		merged = append(merged, v)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Kills > merged[j].Kills
	})

	return merged
}
