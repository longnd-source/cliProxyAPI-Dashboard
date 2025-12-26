package management

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/database"
	log "github.com/sirupsen/logrus"
)

// GetUsageTrends retrieves aggregated usage trends.
// Query params:
// - group_by: "hour" or "day" (default: "hour")
// - limit: number of buckets (default: 24)
func (h *Handler) GetUsageTrends(c *gin.Context) {
	groupBy := c.DefaultQuery("group_by", "hour")
	limitStr := c.DefaultQuery("limit", "24")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 24
	}

	trends, err := database.GetUsageTrends(groupBy, limit)
	if err != nil {
		log.Errorf("failed to fetch usage trends: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch trends"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trends})
}
