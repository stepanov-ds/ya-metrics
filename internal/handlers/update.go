// Package handlers implements HTTP handlers for the metrics server.
//
// It includes:
// - Metric update and retrieval handlers
// - Health check and ping endpoints
// - Root endpoint to list all metrics
package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

// Update handles metric updates via URL parameters or JSON body.
//
// Supports both:
// - URL path format: /update/:metric_type/:metric_name/:value
// - JSON POST format with metric data
//
// Validates input and stores metric in the provided storage.
func Update(c *gin.Context, st storage.Storage) {
	metricType := c.Param("metric_type")
	metricName := c.Param("metric_name")
	metricValue := c.Param("value")
	var m *utils.Metrics
	ctx := c.Request.Context()

	if metricType == "" || metricName == "" || metricValue == "" {
		if c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			if len(body) == 0 {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			m = UpdateWithJSON(c, st)
			metricName = m.ID
			metricType = m.MType
			if strings.ToLower(metricType) == "counter" {
				if m.Delta != nil {
					metricValue = strconv.FormatInt(*m.Delta, 10)
				} else {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
			} else if strings.ToLower(metricType) == "gauge" {
				if m.Value != nil {
					metricValue = strconv.FormatFloat(*m.Value, 'f', -1, 64)
				} else {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
			} else {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}

	if c.Request.Method != http.MethodPost {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	if strings.ToLower(metricType) == "gauge" {
		v, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		st.SetMetric(ctx, metricName, v, strings.ToLower(metricType) == "counter")
	} else if strings.ToLower(metricType) == "counter" {
		v, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		st.SetMetric(ctx, metricName, v, strings.ToLower(metricType) == "counter")
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Data(http.StatusOK, "", nil)
}

// UpdateWithJSON handles metric updates via JSON request body.
//
// Binds incoming JSON to utils.Metrics struct and returns a pointer to it.
// If binding fails, returns an empty Metrics object.
func UpdateWithJSON(c *gin.Context, st storage.Storage) *utils.Metrics {
	var m utils.Metrics
	if err := c.ShouldBindJSON(&m); err != nil {
		return &utils.Metrics{}
	}
	return &m

}
