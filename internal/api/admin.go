package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminStats represents the admin dashboard statistics
type AdminStats struct {
	PlayersReadyToPoll int64 `json:"players_ready_to_poll"`
	PlayersWithErrors  int64 `json:"players_with_errors"`
}

func (s *Server) admin(c *gin.Context) {
	var stats AdminStats
	var err error

	// Count players ready to poll (next_poll_at <= NOW())
	stats.PlayersReadyToPoll, err = s.sqlite.GetPlayersReadyToPollCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count players ready to poll"})
		return
	}

	// Count players with errors (error_count >= 1)
	stats.PlayersWithErrors, err = s.sqlite.GetPlayersWithErrorsCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count players with errors"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
