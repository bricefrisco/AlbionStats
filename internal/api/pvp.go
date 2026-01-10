package api

// import (
// 	"database/sql"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// type playerPvpResponse struct {
// 	Timestamps []time.Time `json:"Timestamps"`
// 	KillFame   []int64     `json:"KillFame"`
// 	DeathFame  []int64     `json:"DeathFame"`
// 	FameRatio  []*float64  `json:"FameRatio"`
// }

// func (s *Server) playerPvp(c *gin.Context) {
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

// 	var playerID string
// 	if err := s.db.
// 		Raw("SELECT player_id FROM player_stats_latest WHERE region = ? AND lower(name) = ?", server, lowerName).
// 		Row().
// 		Scan(&playerID); err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
// 		return
// 	}

// 	rows, err := s.db.
// 		Raw("SELECT ts, kill_fame, death_fame, fame_ratio FROM player_stats_snapshots WHERE region = ? AND player_id = ? ORDER BY ts", server, playerID).
// 		Rows()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
// 		return
// 	}
// 	defer rows.Close()

// 	var (
// 		resp      playerPvpResponse
// 		ts        time.Time
// 		killFame  int64
// 		deathFame int64
// 		fameRatio sql.NullFloat64
// 	)

// 	for rows.Next() {
// 		if err := rows.Scan(&ts, &killFame, &deathFame, &fameRatio); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read player snapshots"})
// 			return
// 		}

// 		resp.Timestamps = append(resp.Timestamps, ts)
// 		resp.KillFame = append(resp.KillFame, killFame)
// 		resp.DeathFame = append(resp.DeathFame, deathFame)

// 		if fameRatio.Valid {
// 			value := fameRatio.Float64
// 			resp.FameRatio = append(resp.FameRatio, &value)
// 		} else {
// 			resp.FameRatio = append(resp.FameRatio, nil)
// 		}
// 	}

// 	if err := rows.Err(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate snapshots"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, resp)
// }
