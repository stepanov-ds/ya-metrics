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

func Update(c *gin.Context, st storage.Storage) {
	metricType := c.Param("metric_type")
	metricName := c.Param("metric_name")
	metricValue := c.Param("value")

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
			UpdateWithJson(c, st)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if c.Request.Method != http.MethodPost {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	switch strings.ToLower(metricType) {
	case "gauge":

		gauge, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		st.SetMetricGauge(metricName, gauge)
	case "counter":
		counter, err := strconv.ParseInt(metricValue, 0, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		st.SetMetricCounter(metricName, counter)
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Data(http.StatusOK, "", nil)
}

func UpdateWithJson(c *gin.Context, st storage.Storage) {
	var m utils.Metrics
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	switch strings.ToLower(m.MType) {
	case "gauge":
		if m.Value == nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		st.SetMetricGauge(m.ID, *m.Value)
	case "counter":
		if m.Delta == nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		st.SetMetricCounter(m.ID, *m.Delta)
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, nil)
}
