package api

import (
	"albionstats/internal/postgres"

	"github.com/gin-gonic/gin"
)

type Server struct {
	postgres *postgres.Postgres
	router   *gin.Engine
}

type Config struct {
	Postgres *postgres.Postgres
}

func NewServer(cfg Config) *Server {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(corsMiddleware())

	server := &Server{
		postgres: cfg.Postgres,
		router:   router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	v1 := s.router.Group("/api")
	// v1.GET("/search/:server/:query", s.search)
	// v1.GET("/metrics/:metricId", s.metrics)
	// v1.GET("/players/:server/:name", s.player)
	// v1.GET("/players/:server/:name/pvp", s.playerPvp)
	// v1.GET("/players/:server/:name/pve", s.playerPve)
	v1.GET("/vm", s.vmQueryProxy)
	v1.GET("/admin", s.admin)
}

func (s *Server) Run(addr string) error {
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
