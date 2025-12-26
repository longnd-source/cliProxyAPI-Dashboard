package management

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for management console
	},
}

// GetEvents upgrades the HTTP connection to a WebSocket and streams real-time usage events.
func (h *Handler) GetEvents(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.WithError(err).Error("failed to upgrade connection for events")
		return
	}
	defer conn.Close()

	if h.usageStats == nil {
		_ = conn.WriteJSON(gin.H{"error": "usage statistics not available"})
		return
	}

	// Subscribe to events
	events, cancel := h.usageStats.Subscribe()
	defer cancel()

	// Handle incoming messages (e.g. ping/close)
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				cancel()
				return
			}
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-events:
			if !ok {
				return
			}
			if err := conn.WriteJSON(event); err != nil {
				if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.WithError(err).Warn("failed to write event to websocket")
				}
				return
			}
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				return
			}
		case <-c.Request.Context().Done():
			return
		}
	}
}


