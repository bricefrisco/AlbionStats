package api

import (
	"albionstats/internal/postgres"
	"albionstats/internal/util"
	"errors"
	"net/http"
	"strconv"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type AllianceOverviewResponse struct {
	Name          string                          `json:"Name"`
	RosterStats   *postgres.PlayerRosterStats     `json:"RosterStats"`
	BattleSummary *postgres.AllianceBattleSummary `json:"BattleSummary"`
	Guilds        []postgres.AllianceGuildStats   `json:"Guilds"`
}

func (s *Server) allianceOverview(c *gin.Context) {
	region := c.Param("server")
	if !util.IsValidServer(region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region"})
		return
	}

	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alliance name is required"})
		return
	}

	player, err := s.postgres.GetPlayerStatsByAllianceName(c.Request.Context(), postgres.Region(region), name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Alliance not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alliance overview"})
		return
	}

	playerCountStr := c.DefaultQuery("playerCount", "10")
	playerCount, err := strconv.Atoi(playerCountStr)
	if err != nil || playerCount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playerCount parameter"})
		return
	}

	var (
		roster  *postgres.PlayerRosterStats
		summary *postgres.AllianceBattleSummary
		guilds  []postgres.AllianceGuildStats
	)

	g, ctx := errgroup.WithContext(c.Request.Context())
	g.Go(func() error {
		var err error
		roster, err = s.postgres.GetAllianceRosterStats(ctx, postgres.Region(region), *player.AllianceName)
		return err
	})
	g.Go(func() error {
		var err error
		summary, err = s.postgres.GetAllianceBattleSummary(ctx, region, *player.AllianceName)
		return err
	})
	g.Go(func() error {
		var err error
		guilds, err = s.postgres.GetAllianceGuildStats(region, *player.AllianceName, playerCount)
		return err
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alliance overview"})
		return
	}

	c.JSON(http.StatusOK, AllianceOverviewResponse{
		Name:          *player.AllianceName,
		RosterStats:   roster,
		BattleSummary: summary,
		Guilds:        guilds,
	})
}
