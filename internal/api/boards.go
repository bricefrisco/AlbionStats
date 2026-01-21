package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) battleSummaries(c *gin.Context) {
	region := c.Param("region")

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	totalPlayersStr := c.DefaultQuery("totalPlayers", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter (must be 1-50)"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	totalPlayers, err := strconv.Atoi(totalPlayersStr)
	if err != nil || totalPlayers < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid totalPlayers parameter"})
		return
	}

	summaries, err := s.postgres.GetBattleSummariesByRegion(region, limit, offset, totalPlayers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get battle summaries"})
		return
	}

	c.JSON(http.StatusOK, summaries)
}

func (s *Server) battleAllianceSummaries(c *gin.Context) {
	region := c.Param("region")
	allianceName := c.Param("allianceName")

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	playerCountStr := c.DefaultQuery("playerCount", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter (must be 1-50)"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	playerCount, err := strconv.Atoi(playerCountStr)
	if err != nil || playerCount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playerCount parameter"})
		return
	}

	summaries, err := s.postgres.GetBattleSummariesByAlliance(region, allianceName, playerCount, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get battle summaries"})
		return
	}

	c.JSON(http.StatusOK, summaries)
}
