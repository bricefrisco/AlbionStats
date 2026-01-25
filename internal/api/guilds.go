package api

import (
	"albionstats/internal/postgres"
	"albionstats/internal/util"
	"errors"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GuildOverviewResponse struct {
	RosterStats   *postgres.PlayerRosterStats  `json:"RosterStats"`
	BattleSummary *postgres.GuildBattleSummary `json:"BattleSummary"`
}

func (s *Server) guildOverview(c *gin.Context) {
	region := c.Param("server")
	if !util.IsValidServer(region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region"})
		return
	}

	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guild name is required"})
		return
	}

	player, err := s.postgres.GetPlayerStatsByGuildName(c.Request.Context(), postgres.Region(region), name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Guild not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find guild"})
		return
	}

	var (
		roster  *postgres.PlayerRosterStats
		summary *postgres.GuildBattleSummary
	)

	g, ctx := errgroup.WithContext(c.Request.Context())
	g.Go(func() error {
		var err error
		roster, err = s.postgres.GetGuildRosterStats(ctx, postgres.Region(region), *player.GuildName)
		return err
	})
	g.Go(func() error {
		var err error
		summary, err = s.postgres.GetGuildBattleSummary(ctx, region, *player.GuildName)
		return err
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get guild overview"})
		return
	}

	c.JSON(http.StatusOK, GuildOverviewResponse{
		RosterStats:   roster,
		BattleSummary: summary,
	})
}
