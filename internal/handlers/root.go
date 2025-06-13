// Package handlers implements HTTP handlers for the metrics server.
//
// It includes:
// - Metric update and retrieval handlers
// - Health check and ping endpoints
// - Root endpoint to list all metrics
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)

// Root handles the root endpoint ("/") and returns all stored metrics in JSON format.
//
// Responds with:
// - 200 OK and JSON body if successful
// - 500 Internal Server Error if JSON marshaling fails
func Root(c *gin.Context, st storage.Storage) {
	jsonData, err := json.Marshal(st.GetAllMetrics())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Writer.Header().Add("Content-Type", "text/html")
	c.String(http.StatusOK, string(jsonData))
}
