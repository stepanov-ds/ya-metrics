package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

func Value(c *gin.Context, st storage.Storage) {
	metricType := c.Param("metric_type")
	metricName := c.Param("metric_name")

	if metricType == "" || metricName == "" {
		if c.Request.Method == http.MethodPost {
			ValueWithJSON(c, st)
			return
		} else {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}

	metricValue, found := st.GetMetric(metricName)
	if found {
		if strings.EqualFold(strings.ToLower(metricType), strings.ToLower(metricValue.MType)) {
			if strings.ToLower(metricValue.MType) == "gauge" {
				result := fmt.Sprintf("%.3f", metricValue.Get())
				result = strings.TrimRight(result, "0")
				result = strings.TrimRight(result, ".")
				c.String(http.StatusOK, result)
				return
			} else {
				c.String(http.StatusOK, fmt.Sprintf("%d", metricValue.Get()))
				return
			}
		} else {
			c.String(http.StatusNotFound, "")
			return
		}
	} else {
		c.String(http.StatusNotFound, "")
		return
	}
}

func ValueWithJSON(c *gin.Context, st storage.Storage) {
	var m utils.Metrics

	if err := c.ShouldBindJSON(&m); err == nil {
		metricValue, found := st.GetMetric(m.ID)
		if found && strings.EqualFold(strings.ToLower(m.MType), strings.ToLower(metricValue.MType)) {
			c.JSON(http.StatusOK, metricValue)
			return
		} else {
			c.String(http.StatusNotFound, "")
			return
		}
	} else {
		println(err.Error())
		c.String(http.StatusBadRequest, "")
		return
	}
}
