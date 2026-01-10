package api

// import (
// 	"database/sql"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// type playerPveResponse struct {
// 	Timestamps   []time.Time `json:"Timestamps"`
// 	PveTotal     []int64     `json:"PveTotal"`
// 	PveRoyal     []int64     `json:"PveRoyal"`
// 	PveOutlands  []int64     `json:"PveOutlands"`
// 	PveAvalon    []int64     `json:"PveAvalon"`
// 	PveHellgate  []int64     `json:"PveHellgate"`
// 	PveCorrupted []int64     `json:"PveCorrupted"`
// 	PveMists     []int64     `json:"PveMists"`
// }

// func (s *Server) playerPve(c *gin.Context) {
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
// 		Raw("SELECT ts, pve_total, pve_royal, pve_outlands, pve_avalon, pve_hellgate, pve_corrupted, pve_mists FROM player_stats_snapshots WHERE region = ? AND player_id = ? ORDER BY ts", server, playerID).
// 		Rows()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
// 		return
// 	}
// 	defer rows.Close()

// 	var (
// 		resp         playerPveResponse
// 		ts           time.Time
// 		pveTotal     int64
// 		pveRoyal     int64
// 		pveOutlands  int64
// 		pveAvalon    int64
// 		pveHellgate  int64
// 		pveCorrupted int64
// 		pveMists     int64
// 	)

// 	for rows.Next() {
// 		if err := rows.Scan(&ts, &pveTotal, &pveRoyal, &pveOutlands, &pveAvalon, &pveHellgate, &pveCorrupted, &pveMists); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read player snapshots"})
// 			return
// 		}

// 		resp.Timestamps = append(resp.Timestamps, ts)
// 		resp.PveTotal = append(resp.PveTotal, pveTotal)
// 		resp.PveRoyal = append(resp.PveRoyal, pveRoyal)
// 		resp.PveOutlands = append(resp.PveOutlands, pveOutlands)
// 		resp.PveAvalon = append(resp.PveAvalon, pveAvalon)
// 		resp.PveHellgate = append(resp.PveHellgate, pveHellgate)
// 		resp.PveCorrupted = append(resp.PveCorrupted, pveCorrupted)
// 		resp.PveMists = append(resp.PveMists, pveMists)
// 	}

// 	if err := rows.Err(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate snapshots"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, resp)
// }
