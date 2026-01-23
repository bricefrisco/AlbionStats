package api

import (
	"albionstats/internal/postgres"
	"albionstats/internal/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlayerStatsResponse struct {
	Player     *postgres.PlayerStatsLatest
	Timestamps []int64
	Pve        *postgres.PlayerPveSeries
	Pvp        *postgres.PlayerPvpSeries
	Gathering  *postgres.PlayerGatheringSeries
	Crafting   *postgres.PlayerCraftingSeries
}

func (s *Server) player(c *gin.Context) {
	server := c.Param("server")
	name := c.Param("name")

	if server == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server is required"})
		return
	}

	if !util.IsValidServer(server) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server. Must be one of: americas, europe, asia"})
		return
	}

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player name is required"})
		return
	}

	region := postgres.Region(server)
	player, err := s.postgres.GetPlayerByName(c.Request.Context(), region, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	statsSeries, err := s.postgres.GetPlayerStatsSeries(region, player.PlayerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	response := PlayerStatsResponse{
		Player:     player,
		Timestamps: statsSeries.Timestamps,
		Pve:        &statsSeries.Pve,
		Pvp:        &statsSeries.Pvp,
		Gathering:  &statsSeries.Gathering,
		Crafting:   &statsSeries.Crafting,
	}

	c.JSON(http.StatusOK, response)
}
