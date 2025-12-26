package management

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/usage"
	coreusage "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/usage"
	"github.com/stretchr/testify/assert"
)

func TestGetEvents_WebSocket(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	usageStats := usage.NewRequestStatistics()
	usage.SetStatisticsEnabled(true)
	
	h := &Handler{
		usageStats: usageStats,
	}
	router.GET("/events", h.GetEvents)

	server := httptest.NewServer(router)
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/events"

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Trigger an event
	go func() {
		time.Sleep(100 * time.Millisecond) // Wait for subscription
		usageStats.Record(context.Background(), coreusage.Record{
			APIKey:      "test-key",
			Model:       "test-model",
			RequestedAt: time.Now(),
			Duration:    50 * time.Millisecond,
			Detail: coreusage.Detail{
				TotalTokens: 100,
			},
		})
	}()

	// Read event
	_ = ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err)

	var event usage.UsageEvent
	err = json.Unmarshal(msg, &event)
	assert.NoError(t, err)

	assert.Equal(t, "test-key", event.APIKey)
	assert.Equal(t, "test-model", event.Model)
	assert.Equal(t, int64(100), event.Tokens.TotalTokens)
	assert.Equal(t, 50*time.Millisecond, event.Duration)
}
