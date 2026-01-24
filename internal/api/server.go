package api

import (
	"albionstats/internal/postgres"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Server struct {
	postgres *postgres.Postgres
	router   *gin.Engine
	topCache *topCache
	logger   *slog.Logger
}

type Config struct {
	Postgres *postgres.Postgres
	Logger   *slog.Logger
}

func NewServer(cfg Config) *Server {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(corsMiddleware())

	server := &Server{
		postgres: cfg.Postgres,
		router:   router,
		topCache: newTopCache(),
		logger:   cfg.Logger,
	}

	server.setupRoutes()
	if err := server.refreshTopCache(); err != nil {
		server.logger.Error("top cache refresh failed", "err", err)
	}
	server.startTopCacheRefresher()
	return server
}

func (s *Server) setupRoutes() {
	s.router.Use(gzip.Gzip(
		gzip.DefaultCompression,
		gzip.WithCustomShouldCompressFn(func(c *gin.Context) bool {
			if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
				return false
			}

			if c.Writer.Header().Get("Content-Encoding") != "" {
				return false
			}

			cl := c.Writer.Header().Get("Content-Length")
			if cl == "" {
				return false
			}

			n, err := strconv.Atoi(cl)
			if err != nil {
				return false
			}

			return n >= 2048
		}),
	))

	// API routes
	v1 := s.router.Group("/api")
	v1.GET("/metrics/admin", s.admin)
	v1.GET("/metrics/:metricId", s.metrics)
	v1.GET("/players/:server/:name", s.player)
	v1.GET("/guilds/:server/:name", s.guildOverview)
	v1.GET("/players/search/:server/:query", s.searchPlayers)
	v1.GET("/guilds/search/:server/:query", s.searchGuilds)
	v1.GET("/alliances/search/:server/:query", s.searchAlliances)
	v1.GET("/alliances/:server/:name", s.allianceOverview)
	v1.GET("/alliances/top/:region", s.topAlliances)
	v1.GET("/guilds/top/:region", s.topGuilds)
	v1.GET("/players/top/:region", s.topPlayers)
	v1.GET("/boards/:region", s.battleSummaries)
	v1.GET("/boards/guild/:region/:guildName", s.battleGuildSummaries)
	v1.GET("/boards/alliance/:region/:allianceName", s.battleAllianceSummaries)
	v1.GET("/boards/player/:region/:playerName", s.battlePlayerSummaries)
	v1.GET("/battles/:region/:battleId", s.battle)
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
