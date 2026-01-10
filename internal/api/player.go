// package api

// import (
// 	"albionstats/internal/database"
// 	"errors"
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// func (s *Server) player(c *gin.Context) {
// 	server := c.Param("server")
// 	name := c.Param("name")

// 	if server == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Server is required"})
// 		return
// 	}

// 	validServers := map[string]bool{
// 		"americas": true,
// 		"europe":   true,
// 		"asia":     true,
// 	}

// 	if !validServers[server] {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server. Must be one of: americas, europe, asia"})
// 		return
// 	}

// 	if name == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Player name is required"})
// 		return
// 	}

// 	lowerName := strings.ToLower(name)

// 	var player database.PlayerStatsLatest
// 	err := s.db.
// 		Where("region = ? AND LOWER(name) = ?", server, lowerName).
// 		First(&player).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, player)
// }
