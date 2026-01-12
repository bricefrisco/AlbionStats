package api

import (
	"albionstats/internal/postgres"

	"github.com/gin-contrib/gzip"
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
	s.router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Static front-end files
	s.router.Static("/_app", "/usr/local/share/albionstats_web/_app")
	s.router.StaticFile("/robots.txt", "/usr/local/share/albionstats_web/robots.txt")
	s.router.NoRoute(func(c *gin.Context) {
		c.File("/usr/local/share/albionstats_web/200.html")
	})

	// API routes
	v1 := s.router.Group("/api")
	v1.GET("/metrics/admin", s.admin)
	v1.GET("/metrics/:metricId", s.metrics)
	v1.GET("/metrics/pvp/:server/:playerId", s.playerPvp)
	v1.GET("/metrics/pve/:server/:playerId", s.playerPve)
	v1.GET("/metrics/gathering/:server/:playerId", s.playerGathering)
	v1.GET("/metrics/crafting/:server/:playerId", s.playerCrafting)
	v1.GET("/players/:server/:name", s.player)
	v1.GET("/search/:server/:query", s.search)
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
