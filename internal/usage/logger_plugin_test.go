package usage

import (
	"context"
	"testing"
	"time"

	coreusage "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/usage"
	"github.com/stretchr/testify/assert"
)

func TestLoggerPlugin_LatencyAndLimits(t *testing.T) {
	// Setup
	p := NewLoggerPlugin()
	ctx := context.Background()
	apiKey := "test-api-key"
	model := "test-model"

	// Test 1: Verify Duration recording
	start := time.Now()
	duration := 100 * time.Millisecond
	record := coreusage.Record{
		APIKey:      apiKey,
		Model:       model,
		RequestedAt: start.Add(-duration),
		Duration:    duration,
		Detail: coreusage.Detail{
			TotalTokens: 10,
		},
	}
	p.HandleUsage(ctx, record)

	snapshot := p.stats.Snapshot()
	assert.Equal(t, int64(1), snapshot.TotalRequests)
	
	apiStats, ok := snapshot.APIs[apiKey]
	assert.True(t, ok)
	
	modelStats, ok := apiStats.Models[model]
	assert.True(t, ok)
	assert.Len(t, modelStats.Details, 1)
	assert.Equal(t, duration, modelStats.Details[0].Duration)

	// Test 2: Verify Memory Limit (Cap at 100)
	// We already added 1, add 105 more
	for i := 0; i < 105; i++ {
		p.HandleUsage(ctx, record)
	}

	snapshot = p.stats.Snapshot()
	apiStats = snapshot.APIs[apiKey]
	modelStats = apiStats.Models[model]

	assert.Equal(t, int64(106), modelStats.TotalRequests) // 1 + 105
	assert.Len(t, modelStats.Details, 100)                // Capped at 100
}
