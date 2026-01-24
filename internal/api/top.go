package api

import (
	"albionstats/internal/postgres"
	"albionstats/internal/util"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const topCacheRefreshInterval = 5 * time.Minute

type topCache struct {
	mu        sync.RWMutex
	alliances map[string][]postgres.TopAllianceStats
	guilds    map[string][]postgres.TopGuildStats
	players   map[string][]postgres.TopPlayerStats
}

func newTopCache() *topCache {
	return &topCache{
		alliances: make(map[string][]postgres.TopAllianceStats),
		guilds:    make(map[string][]postgres.TopGuildStats),
		players:   make(map[string][]postgres.TopPlayerStats),
	}
}

func (s *Server) refreshTopCache() error {
	if s.topCache == nil {
		return nil
	}

	regions := []string{"americas", "europe", "asia"}
	for _, region := range regions {
		alliances, err := s.postgres.GetTopAlliances(region, 50, 0)
		if err != nil {
			return err
		}

		guilds, err := s.postgres.GetTopGuilds(region, 50, 0)
		if err != nil {
			return err
		}

		players, err := s.postgres.GetTopPlayers(region, 50, 0)
		if err != nil {
			return err
		}

		s.topCache.set(region, alliances, guilds, players)
	}
	return nil
}

func (s *Server) startTopCacheRefresher() {
	ticker := time.NewTicker(topCacheRefreshInterval)
	go func() {
		for range ticker.C {
			err := s.refreshTopCache()
			if err != nil {
				s.logger.Error("top cache refresh failed", "err", err)
			}
		}
	}()
}

func (c *topCache) set(region string, alliances []postgres.TopAllianceStats, guilds []postgres.TopGuildStats, players []postgres.TopPlayerStats) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.alliances[region] = append([]postgres.TopAllianceStats(nil), alliances...)
	c.guilds[region] = append([]postgres.TopGuildStats(nil), guilds...)
	c.players[region] = append([]postgres.TopPlayerStats(nil), players...)
}

func (c *topCache) getAlliances(region string) ([]postgres.TopAllianceStats, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats, ok := c.alliances[region]
	return stats, ok
}

func (c *topCache) getGuilds(region string) ([]postgres.TopGuildStats, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats, ok := c.guilds[region]
	return stats, ok
}

func (c *topCache) getPlayers(region string) ([]postgres.TopPlayerStats, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats, ok := c.players[region]
	return stats, ok
}

func (s *Server) topAlliances(c *gin.Context) {
	region := c.Param("region")
	if !util.IsValidServer(region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region"})
		return
	}

	limit, offset, ok := parseLimitOffset(c)
	if !ok {
		return
	}

	if s.topCache == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Top cache is not initialized"})
		return
	}

	stats, ok := s.topCache.getAlliances(region)
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Top alliances not ready"})
		return
	}

	c.JSON(http.StatusOK, sliceTopAlliances(stats, limit, offset))
}

func (s *Server) topGuilds(c *gin.Context) {
	region := c.Param("region")
	if !util.IsValidServer(region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region"})
		return
	}

	limit, offset, ok := parseLimitOffset(c)
	if !ok {
		return
	}

	if s.topCache == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Top cache is not initialized"})
		return
	}

	stats, ok := s.topCache.getGuilds(region)
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Top guilds not ready"})
		return
	}

	c.JSON(http.StatusOK, sliceTopGuilds(stats, limit, offset))
}

func (s *Server) topPlayers(c *gin.Context) {
	region := c.Param("region")
	if !util.IsValidServer(region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region"})
		return
	}

	limit, offset, ok := parseLimitOffset(c)
	if !ok {
		return
	}

	if s.topCache == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Top cache is not initialized"})
		return
	}

	stats, ok := s.topCache.getPlayers(region)
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Top players not ready"})
		return
	}

	c.JSON(http.StatusOK, sliceTopPlayers(stats, limit, offset))
}

func parseLimitOffset(c *gin.Context) (int, int, bool) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter (must be 1-50)"})
		return 0, 0, false
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return 0, 0, false
	}

	return limit, offset, true
}

func sliceTopAlliances(stats []postgres.TopAllianceStats, limit int, offset int) []postgres.TopAllianceStats {
	if offset >= len(stats) {
		return []postgres.TopAllianceStats{}
	}
	end := offset + limit
	if end > len(stats) {
		end = len(stats)
	}
	return stats[offset:end]
}

func sliceTopGuilds(stats []postgres.TopGuildStats, limit int, offset int) []postgres.TopGuildStats {
	if offset >= len(stats) {
		return []postgres.TopGuildStats{}
	}
	end := offset + limit
	if end > len(stats) {
		end = len(stats)
	}
	return stats[offset:end]
}

func sliceTopPlayers(stats []postgres.TopPlayerStats, limit int, offset int) []postgres.TopPlayerStats {
	if offset >= len(stats) {
		return []postgres.TopPlayerStats{}
	}
	end := offset + limit
	if end > len(stats) {
		end = len(stats)
	}
	return stats[offset:end]
}
