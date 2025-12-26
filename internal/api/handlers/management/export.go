package management

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/database"
)

// ExportActivityLogs generates and downloads a CSV of the activity logs.
func (h *Handler) ExportActivityLogs(c *gin.Context) {
	modelFilter := c.Query("model")
	statusFilter := c.Query("status")

	logs, err := database.GetAllActivity(modelFilter, statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch logs for export: %v", err)})
		return
	}

	// Set headers for file download
	filename := fmt.Sprintf("activity_logs_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	
	// Write CSV Header
	header := []string{"ID", "Timestamp", "Model", "Source", "Input Tokens", "Output Tokens", "Total Tokens", "Duration (ms)", "Status", "API Key"}
	if err := writer.Write(header); err != nil {
		return // Streaming error, can't json response now
	}

	// Write Rows
	for _, log := range logs {
		status := "Success"
		if log.IsFailure {
			status = "Failure"
		}
		row := []string{
			strconv.FormatInt(log.ID, 10),
			log.Timestamp.Format(time.RFC3339),
			log.Model,
			log.Source,
			strconv.FormatInt(log.InputTokens, 10),
			strconv.FormatInt(log.OutputTokens, 10),
			strconv.FormatInt(log.TotalTokens, 10),
			strconv.FormatInt(log.DurationMs, 10),
			status,
			log.APIKey, // Note: You might want to mask this if sensitive, but this is admin export
		}
		if err := writer.Write(row); err != nil {
			return
		}
	}

	writer.Flush()
}
