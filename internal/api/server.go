package api

import (
	"albionstats/internal/models"
	"net/http"

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
	v1 := s.router.Group("/v1")
	v1.GET("/search/:query", s.search)
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
	err := s.db.Select("player_id", "name").
		Where("LOWER(name) LIKE LOWER(?)", query+"%").
		Limit(6).
		Order("name ASC").
		Find(&players).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	type PlayerSearchResult struct {
		PlayerID string `json:"player_id"`
		Name     string `json:"name"`
	}

	results := make([]PlayerSearchResult, len(players))
	for i, player := range players {
		results[i] = PlayerSearchResult{
			PlayerID: player.PlayerID,
			Name:     player.Name,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"players": results,
	})
}
