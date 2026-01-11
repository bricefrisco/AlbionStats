package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) metrics(c *gin.Context) {
	metricId := c.Param("metricId")
	if metricId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric ID is required"})
		return
	}

	granularity := c.DefaultQuery("granularity", "1w")

	results, err := s.postgres.GetMetrics(c.Request.Context(), metricId, granularity)
	if err != nil {
		// Check if it's an invalid granularity error
		if err.Error() == "invalid granularity: "+granularity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid granularity. Use: 1w, 1m, 1y, or all"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}
