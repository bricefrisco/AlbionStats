package api

import (
	"albionstats/internal/postgres"
	"albionstats/internal/util"
	"errors"
	"net/http"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type AllianceOverviewResponse struct {
	Name          string                          `json:"Name"`
	RosterStats   *postgres.PlayerRosterStats     `json:"RosterStats"`
	BattleSummary *postgres.AllianceBattleSummary `json:"BattleSummary"`
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
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alliance overview"})
	}

	var (
		roster  *postgres.PlayerRosterStats
		summary *postgres.AllianceBattleSummary
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

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alliance overview"})
		return
	}

	c.JSON(http.StatusOK, AllianceOverviewResponse{
		Name:          *player.AllianceName,
		RosterStats:   roster,
		BattleSummary: summary,
	})
}
