package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)

func Value(c *gin.Context, st storage.Storage) {
	metricType := c.Param("metric_type")
	metricName := c.Param("metric_name")

	if metricType == "" || metricName == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	metricValue, found := st.GetMetric(metricName)
	if found {
		if metricValue.IsCounter && strings.ToLower(metricType) == "counter" {
			c.String(http.StatusOK, fmt.Sprintf("%d", metricValue.Counter))
		} else if !metricValue.IsCounter && strings.ToLower(metricType) == "gauge" {
			result := fmt.Sprintf("%.3f", metricValue.Gauge)
			result = strings.TrimRight(result, "0")
			result = strings.TrimRight(result, ".")
			c.String(http.StatusOK, result)
		} else {
			c.String(http.StatusNotFound, "")
		}
	} else {
		c.String(http.StatusNotFound, "")
	}
}
