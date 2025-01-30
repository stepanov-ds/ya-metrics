package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/cmd/server/storage"
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
		if metricValue.IsCounter {
			c.String(http.StatusOK, fmt.Sprintf("%d", metricValue.Counter))
		} else {
			c.String(http.StatusOK, fmt.Sprintf("%f", metricValue.Gauge))
		}
	} else {
		c.String(http.StatusNotFound, "")
	}
}
