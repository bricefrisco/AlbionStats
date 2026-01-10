package api

// import (
// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// type Server struct {
// 	db     *gorm.DB
// 	router *gin.Engine
// }

// type Config struct {
// 	Port string
// }

// func New(db *gorm.DB, cfg Config) *Server {
// 	gin.SetMode(gin.ReleaseMode) // production mode

// 	router := gin.Default()
// 	router.Use(corsMiddleware())

// 	server := &Server{
// 		db:     db,
// 		router: router,
// 	}

// 	server.setupRoutes()
// 	return server
// }

// func (s *Server) setupRoutes() {
// 	v1 := s.router.Group("/albionstats/v1")
// 	v1.GET("/search/:server/:query", s.search)
// 	v1.GET("/metrics/:metricId", s.metrics)
// 	v1.GET("/players/:server/:name", s.player)
// 	v1.GET("/players/:server/:name/pvp", s.playerPvp)
// 	v1.GET("/players/:server/:name/pve", s.playerPve)
// }

// func (s *Server) Run(addr string) error {
// 	return s.router.Run(addr)
// }

// func (s *Server) Router() *gin.Engine {
// 	return s.router
// }

// // CORS middleware
// func corsMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Credentials", "true")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }
