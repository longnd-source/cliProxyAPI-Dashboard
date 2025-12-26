package database

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDatabase_Flow(t *testing.T) {
	// Setup temp dir
	tmpDir, err := os.MkdirTemp("", "cliproxy-db-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Init
	err = Init(tmpDir)
	require.NoError(t, err)
	defer Close()

	// Insert Log
	now := time.Now()
	log := UsageLog{
		Timestamp:   now,
		APIKey:      "test-key",
		Model:       "test-model",
		InputTokens: 10,
		OutputTokens: 20,
		TotalTokens: 30,
		IsFailure:   false,
		Source:      "test",
	}
	err = InsertUsageLog(log)
	require.NoError(t, err)

	// Get Recent Activity
	logs, err := GetRecentActivity(10, 0, "", "")
	require.NoError(t, err)
	require.Len(t, logs, 1)
	t.Logf("Retrieved Log Timestamp: %v", logs[0].Timestamp)

	// Get Trends
	trends, err := GetUsageTrends("hour", 24)
	require.NoError(t, err)
	require.NotEmpty(t, trends, "Trends should not be empty")
	require.Equal(t, int64(1), trends[len(trends)-1].Requests)
}
