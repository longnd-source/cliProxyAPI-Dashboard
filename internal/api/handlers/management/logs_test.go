package management

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLogs_Search(t *testing.T) {
	// Setup temporary log directory
	tmpDir, err := os.MkdirTemp("", "logs_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a dummy log file
	logContent := `2024-01-01 10:00:00 [INFO] System started
2024-01-01 10:00:01 [ERROR] Database connection failed
2024-01-01 10:00:02 [INFO] Retrying connection
2024-01-01 10:00:03 [ERROR] Fatal error occurred
2024-01-01 10:00:04 [DEBUG] Debugging info
`
	err = os.WriteFile(filepath.Join(tmpDir, "main.log"), []byte(logContent), 0644)
	require.NoError(t, err)

	// Initialize handler
	cfg := &config.Config{
		LoggingToFile: true,
	}
	h := NewHandler(cfg, "", nil)
	h.SetLogDirectory(tmpDir)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/logs", h.GetLogs)

	t.Run("Search for ERROR", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/logs?search=error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Lines []string `json:"lines"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Len(t, resp.Lines, 2)
		assert.Contains(t, resp.Lines[0], "Database connection failed")
		assert.Contains(t, resp.Lines[1], "Fatal error occurred")
	})

	t.Run("Search for info (case insensitive)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/logs?search=INFO", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Lines []string `json:"lines"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Len(t, resp.Lines, 3)
		assert.Contains(t, resp.Lines[0], "System started")
		assert.Contains(t, resp.Lines[1], "Retrying connection")
		assert.Contains(t, resp.Lines[2], "Debugging info")
	})

	t.Run("No search term", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/logs", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Lines []string `json:"lines"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Len(t, resp.Lines, 5)
	})
    
    t.Run("Search not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/logs?search=nonexistent", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Lines []string `json:"lines"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Len(t, resp.Lines, 0)
	})
}
