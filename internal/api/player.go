package api

import (
	"albionstats/internal/postgres"
	"albionstats/internal/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type PlayerStatsResponse struct {
	Player    *postgres.PlayerStatsLatest    `json:"player"`
	Pve       *postgres.PlayerPveStats       `json:"pve"`
	Pvp       *postgres.PlayerPvpStats       `json:"pvp"`
	Gathering *postgres.PlayerGatheringStats `json:"gathering"`
	Crafting  *postgres.PlayerCraftingStats  `json:"crafting"`
}

func (s *Server) player(c *gin.Context) {
	server := c.Param("server")
	name := c.Param("name")

	if server == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server is required"})
		return
	}

	if !util.IsValidServer(server) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server. Must be one of: americas, europe, asia"})
		return
	}

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player name is required"})
		return
	}

	region := postgres.Region(server)
	player, err := s.postgres.GetPlayerByName(c.Request.Context(), region, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	var (
		pveStats       *postgres.PlayerPveStats
		pvpStats       *postgres.PlayerPvpStats
		gatheringStats *postgres.PlayerGatheringStats
		craftingStats  *postgres.PlayerCraftingStats
	)

	g, _ := errgroup.WithContext(c.Request.Context())

	g.Go(func() error {
		var err error
		pveStats, err = s.postgres.GetPlayerPveStats(region, player.PlayerID)
		return err
	})

	g.Go(func() error {
		var err error
		pvpStats, err = s.postgres.GetPlayerPvpStats(region, player.PlayerID)
		return err
	})

	g.Go(func() error {
		var err error
		gatheringStats, err = s.postgres.GetPlayerGatheringStats(region, player.PlayerID)
		return err
	})

	g.Go(func() error {
		var err error
		craftingStats, err = s.postgres.GetPlayerCraftingStats(region, player.PlayerID)
		return err
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	response := PlayerStatsResponse{
		Player:    player,
		Pve:       pveStats,
		Pvp:       pvpStats,
		Gathering: gatheringStats,
		Crafting:  craftingStats,
	}

	c.JSON(http.StatusOK, response)
}
