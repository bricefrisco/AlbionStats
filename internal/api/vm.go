package api

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) vmQueryProxy(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query"})
		return
	}

	// Hard safety gate
	if !isAllowedQuery(query) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "query not allowed",
		})
		return
	}

	// Build upstream request
	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8428/api/v1/query_range",
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "request build failed"})
		return
	}

	q := req.URL.Query()
	for k, v := range c.Request.URL.Query() {
		q[k] = v
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "vm unavailable"})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

var allowedMetric = regexp.MustCompile(`\balbion_player_[a-zA-Z0-9_]*\b`)

func isAllowedQuery(q string) bool {
	// Must reference at least one allowed metric
	if !allowedMetric.MatchString(q) {
		return false
	}

	// Block label introspection / meta queries
	blocked := []string{
		"__name__",
		"label_values",
		"label_names",
		"metrics",
	}

	for _, b := range blocked {
		if strings.Contains(q, b) {
			return false
		}
	}

	return true
}
