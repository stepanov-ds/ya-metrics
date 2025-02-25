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
	var m utils.Metrics

	if metricType == "" || metricName == "" {
		if c.Request.Method == http.MethodPost {
			m = ValueWithJson(c, st)
			metricName = m.ID
			metricType = m.MType
			if metricType == "" || metricName == "" {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
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
			} else {
				c.String(http.StatusOK, fmt.Sprintf("%d", metricValue.Get()))
			}
		} else {
			c.String(http.StatusNotFound, "")
		}
	} else {
		c.String(http.StatusNotFound, "")
	}
}

func ValueWithJson(c *gin.Context, st storage.Storage) utils.Metrics {
	var m utils.Metrics

	if err := c.ShouldBindJSON(&m); err == nil {
		return m
	} else {
		println(err.Error())
		return utils.Metrics{}
	}
}