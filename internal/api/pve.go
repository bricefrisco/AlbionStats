package api

import (
	"net/http"

	"albionstats/internal/postgres"
	"albionstats/internal/util"

	"github.com/gin-gonic/gin"
)

func (s *Server) playerPve(c *gin.Context) {
	server := c.Param("server")
	playerID := c.Param("playerId")

	if server == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server is required"})
		return
	}

	if !util.IsValidServer(server) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server. Must be one of: americas, europe, asia"})
		return
	}

	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player ID is required"})
		return
	}

	stats, err := s.postgres.GetPlayerPveStats(postgres.Region(server), playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch PvE stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
