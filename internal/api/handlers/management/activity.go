package management

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/database"
	log "github.com/sirupsen/logrus"
)

// GetActivityLogs returns paginated activity logs from the database.
func (h *Handler) GetActivityLogs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Cap limit
	if limit > 100 {
		limit = 100
	}

	modelFilter := c.Query("model")
	statusFilter := c.Query("status")

	logs, err := database.GetRecentActivity(limit, offset, modelFilter, statusFilter)
	if err != nil {
		log.Errorf("failed to fetch activity logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch logs: %v", err)})
		return
	}
	if logs == nil {
		logs = []database.UsageLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"meta": gin.H{
			"limit":  limit,
			"offset": offset,
		},
	})
}
