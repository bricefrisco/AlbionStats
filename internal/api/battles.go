package api

import (
	"albionstats/internal/postgres"
	"albionstats/internal/util"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type BattleDetailResponse struct {
	Summary       postgres.BattleSummary         `json:"summary"`
	AllianceStats []postgres.BattleAllianceStats `json:"alliance_stats"`
	GuildStats    []postgres.BattleGuildStats    `json:"guild_stats"`
	PlayerStats   []postgres.BattlePlayerStats   `json:"player_stats"`
}

func (s *Server) battle(c *gin.Context) {
	region := c.Param("region")
	if !util.IsValidServer(region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region"})
		return
	}

	battleIDStr := c.Param("battleId")

	battleID, err := strconv.ParseInt(battleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid battleId"})
		return
	}

	var resp BattleDetailResponse
	var summary postgres.BattleSummary
	var allianceStats []postgres.BattleAllianceStats
	var guildStats []postgres.BattleGuildStats
	var playerStats []postgres.BattlePlayerStats

	var wg sync.WaitGroup
	var errSummary, errAlliance, errGuild, errPlayer error

	wg.Add(4)

	go func() {
		defer wg.Done()
		summary, errSummary = s.postgres.GetBattleSummary(region, battleID)
	}()

	go func() {
		defer wg.Done()
		allianceStats, errAlliance = s.postgres.GetBattleAllianceStats(region, battleID)
	}()

	go func() {
		defer wg.Done()
		guildStats, errGuild = s.postgres.GetBattleGuildStats(region, battleID)
	}()

	go func() {
		defer wg.Done()
		playerStats, errPlayer = s.postgres.GetBattlePlayerStats(region, battleID)
	}()

	wg.Wait()

	if errSummary != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Battle not found"})
		return
	}
	if errAlliance != nil || errGuild != nil || errPlayer != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch battle details"})
		return
	}

	resp.Summary = summary
	resp.AllianceStats = allianceStats
	resp.GuildStats = guildStats
	resp.PlayerStats = playerStats

	c.JSON(http.StatusOK, resp)
}
