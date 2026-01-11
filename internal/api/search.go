package api

import (
	"net/http"

	"albionstats/internal/postgres"

	"github.com/gin-gonic/gin"
)

func (s *Server) search(c *gin.Context) {
	server := c.Param("server")
	query := c.Param("query")

	if server == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server is required"})
		return
	}

	// Validate server parameter
	validServers := map[string]bool{
		"americas": true,
		"europe":   true,
		"asia":     true,
	}
	if !validServers[server] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server. Must be one of: americas, europe, asia"})
		return
	}

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query is required"})
		return
	}

	players, err := s.postgres.SearchPlayers(c.Request.Context(), postgres.Region(server), query, 6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	type PlayerSearchResult struct {
		PlayerID     string  `json:"player_id"`
		Name         string  `json:"name"`
		GuildName    *string `json:"guild_name,omitempty"`
		AllianceName *string `json:"alliance_name,omitempty"`
	}

	results := make([]PlayerSearchResult, len(players))
	for i, player := range players {
		results[i] = PlayerSearchResult{
			PlayerID:     player.PlayerID,
			Name:         player.Name,
			GuildName:    player.GuildName,
			AllianceName: player.AllianceName,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"players": results,
	})
}
