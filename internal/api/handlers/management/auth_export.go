package management

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ExportAllAuthFiles packages all .json files in the auth-dir into a ZIP.
func (h *Handler) ExportAllAuthFiles(c *gin.Context) {
	if h.cfg.AuthDir == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auth-dir not configured"})
		return
	}

	entries, err := os.ReadDir(h.cfg.AuthDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read auth directory: %v", err)})
		return
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	count := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			continue
		}

		f, err := os.Open(filepath.Join(h.cfg.AuthDir, entry.Name()))
		if err != nil {
			continue
		}

		w, err := zw.Create(entry.Name())
		if err != nil {
			f.Close()
			continue
		}

		if _, err := io.Copy(w, f); err != nil {
			f.Close()
			continue
		}
		f.Close()
		count++
	}

	if err := zw.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to finalize ZIP: %v", err)})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no accounts found to export"})
		return
	}

	filename := fmt.Sprintf("cliproxy_accounts_%s.zip", time.Now().Format("20060102_150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/zip")
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}

// ImportAllAuthFiles extracts .json files from a ZIP and registers them.
func (h *Handler) ImportAllAuthFiles(c *gin.Context) {
	if h.authManager == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "core auth manager unavailable"})
		return
	}

	file, err := c.FormFile("bundle")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'bundle' file (ZIP)"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to open uploaded file: %v", err)})
		return
	}
	defer f.Close()

	// Read into memory or use LimitReader to prevent massive file swap
	// For simplicity and since these are small JSONs, let's use a temporary file or memory
	// ZIP reader needs ReaderAt, so we either read to memory or use a temp file.
	// We'll read to memory for now as auth files are tiny.
	body, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read bundle"})
		return
	}

	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ZIP archive"})
		return
	}

	imported := 0
	errors := []string{}
	ctx := c.Request.Context()

	for _, zf := range zr.File {
		if zf.FileInfo().IsDir() || !strings.HasSuffix(strings.ToLower(zf.Name), ".json") {
			continue
		}

		// Sanitize filename to prevent path traversal
		name := filepath.Base(zf.Name)
		dst := filepath.Join(h.cfg.AuthDir, name)

		rc, err := zf.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to open", name))
			continue
		}

		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to read", name))
			continue
		}

		if err := os.WriteFile(dst, data, 0600); err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to save", name))
			continue
		}

		if err := h.registerAuthFromFile(ctx, dst, data); err != nil {
			errors = append(errors, fmt.Sprintf("%s: registration failed: %v", name, err))
			continue
		}

		imported++
	}

	if imported == 0 && len(errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "errors": errors})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"imported": imported,
		"errors":   errors,
	})
}
