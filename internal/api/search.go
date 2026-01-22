package api

import (
	"net/http"

	"albionstats/internal/postgres"
	"albionstats/internal/util"

	"github.com/gin-gonic/gin"
)

func (s *Server) searchPlayers(c *gin.Context) {
	server := c.Param("server")
	query := c.Param("query")

	if server == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server is required"})
		return
	}

	if !util.IsValidServer(server) {
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

func (s *Server) searchAlliances(c *gin.Context) {
	server := c.Param("server")
	query := c.Param("query")

	if server == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server is required"})
		return
	}

	if !util.IsValidServer(server) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server. Must be one of: americas, europe, asia"})
		return
	}

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query is required"})
		return
	}

	alliances, err := s.postgres.SearchAlliances(c.Request.Context(), postgres.Region(server), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alliances": alliances,
	})
}

func (s *Server) searchGuilds(c *gin.Context) {
	server := c.Param("server")
	query := c.Param("query")

	if server == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server is required"})
		return
	}

	if !util.IsValidServer(server) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server. Must be one of: americas, europe, asia"})
		return
	}

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query is required"})
		return
	}

	guilds, err := s.postgres.SearchGuilds(c.Request.Context(), postgres.Region(server), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"guilds": guilds,
	})
}
