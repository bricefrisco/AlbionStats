package api

import (
	"albionstats/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db     *gorm.DB
	router *gin.Engine
}

type Config struct {
	Port string
}

func New(db *gorm.DB, cfg Config) *Server {
	gin.SetMode(gin.ReleaseMode) // production mode

	router := gin.Default()
	router.Use(corsMiddleware())

	server := &Server{
		db:     db,
		router: router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	v1 := s.router.Group("/albionstats/v1")
	v1.GET("/search/:query", s.search)
	v1.GET("/metrics/:metricId", s.metrics)
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) search(c *gin.Context) {
	query := c.Param("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query is required"})
		return
	}

	var players []models.PlayerState
	err := s.db.Select("player_id", "name", "guild_name", "alliance_name").
		Where("LOWER(name) LIKE LOWER(?)", query+"%").
		Limit(6).
		Order("name ASC").
		Find(&players).Error

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

func (s *Server) metrics(c *gin.Context) {
	metricId := c.Param("metricId")
	if metricId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric ID is required"})
		return
	}

	granularity := c.DefaultQuery("granularity", "1w")

	type MetricPoint struct {
		Timestamp time.Time `json:"timestamp"`
		Value     int64     `json:"value"`
	}

	var query string
	var args []interface{}

	switch granularity {
	case "1w":
		// Aggregate to hourly buckets for 1 week
		query = `
			SELECT
				time_bucket('1 hour', ts) AS timestamp,
				max(value) AS value
			FROM metrics_timeseries
			WHERE metric = $1 AND ts >= NOW() - INTERVAL '1 week'
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	case "1m":
		// Aggregate to daily buckets for 1 month
		query = `
			SELECT
				time_bucket('1 day', ts) AS timestamp,
				max(value) AS value
			FROM metrics_timeseries
			WHERE metric = $1 AND ts >= NOW() - INTERVAL '1 month'
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	case "1y":
		// Aggregate to weekly buckets for 1 year
		query = `
			SELECT
				time_bucket('1 week', ts) AS timestamp,
				max(value) AS value
			FROM metrics_timeseries
			WHERE metric = $1 AND ts >= NOW() - INTERVAL '1 year'
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	case "all":
		// Aggregate to monthly buckets for all time
		query = `
			SELECT
				time_bucket('1 month', ts) AS timestamp,
				max(value) AS value
			FROM metrics_timeseries
			WHERE metric = $1
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid granularity. Use: 1w, 1m, 1y, or all"})
		return
	}

	rows, err := s.db.Raw(query, args...).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	var results []MetricPoint
	for rows.Next() {
		var timestamp time.Time
		var value float64
		if err := rows.Scan(&timestamp, &value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		results = append(results, MetricPoint{
			Timestamp: timestamp,
			Value:     int64(value),
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Row iteration error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}
