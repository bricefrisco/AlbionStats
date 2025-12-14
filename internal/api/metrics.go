package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) metrics(c *gin.Context) {
	metricId := c.Param("metricId")
	if metricId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric ID is required"})
		return
	}

	granularity := c.DefaultQuery("granularity", "1w")

	type MetricPoint struct {
		Timestamp time.Time `json:"timestamp"`
		Value     int64     `json:"value"`
	}

	var query string
	var args []interface{}

	switch granularity {
	case "1w":
		// Aggregate to hourly buckets for 1 week
		query = `
			SELECT
				time_bucket('1 hour', ts) AS timestamp,
				max(value) AS value
			FROM metrics_timeseries
			WHERE metric = $1 AND ts >= NOW() - INTERVAL '1 week'
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	case "1m":
		// Aggregate to daily buckets for 1 month
		query = `
			SELECT
				time_bucket('1 day', ts) AS timestamp,
				max(value) AS value
			FROM metrics_timeseries
			WHERE metric = $1 AND ts >= NOW() - INTERVAL '1 month'
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	case "1y":
		// Aggregate to weekly buckets for 1 year
		query = `
			SELECT
				time_bucket('1 week', ts) AS timestamp,
				max(value) AS value
			FROM metrics_timeseries
			WHERE metric = $1 AND ts >= NOW() - INTERVAL '1 year'
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	case "all":
		// Aggregate to monthly buckets for all time
		query = `
			SELECT
				time_bucket('1 month', ts) AS timestamp,
				max(value) AS value
			FROM metrics_timeseries
			WHERE metric = $1
			GROUP BY 1
			ORDER BY 1`
		args = []interface{}{metricId}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid granularity. Use: 1w, 1m, 1y, or all"})
		return
	}

	rows, err := s.db.Raw(query, args...).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	var results []MetricPoint
	for rows.Next() {
		var timestamp time.Time
		var value float64
		if err := rows.Scan(&timestamp, &value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		results = append(results, MetricPoint{
			Timestamp: timestamp,
			Value:     int64(value),
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Row iteration error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}
