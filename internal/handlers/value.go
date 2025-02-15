package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

func Value(c *gin.Context, st storage.Storage) {
	metricType := c.Param("metric_type")
	metricName := c.Param("metric_name")

	if metricType == "" || metricName == "" {
		ValueWithJson(c, st)
		return
	}
	metricValue, found := st.GetMetric(metricName)
	if found {
		switch v := metricValue.(type) {
		case *utils.MetricCounter:
			if strings.ToLower(metricType) == "counter" {
				c.String(http.StatusOK, fmt.Sprintf("%v", v.Get()))
			} else {
				c.String(http.StatusNotFound, "")
			}
		case *utils.MetricGauge:
			if strings.ToLower(metricType) == "gauge" {
				result := fmt.Sprintf("%.3f", v.Get())
				result = strings.TrimRight(result, "0")
				result = strings.TrimRight(result, ".")
				c.String(http.StatusOK, result)
			} else {
				c.String(http.StatusNotFound, "")
			}
		}
	} else {
		c.String(http.StatusNotFound, "")
	}
}

func ValueWithJson(c *gin.Context, st storage.Storage) {
	var m utils.Metrics
	
	defer c.Request.Body.Close()
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if m.ID == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	switch strings.ToLower(m.MType) {
	case "gauge":
		metric, found := st.GetMetric(m.ID)
		if !found {
			c.JSON(http.StatusNotFound, nil)
		} 
		if reflect.TypeOf(metric) != reflect.TypeOf(&utils.MetricGauge{}) {
			c.JSON(http.StatusNotFound, nil)
		}
		floatValue, ok := metric.Get().(float64) 
		if !ok {
			c.JSON(http.StatusNotFound, nil)
		}
		m.Value = &floatValue
		c.JSON(http.StatusOK, m)
	case "counter":
		metric, found := st.GetMetric(m.ID)
		if !found {
			c.JSON(http.StatusNotFound, nil)
		} 
		if reflect.TypeOf(metric) != reflect.TypeOf(&utils.MetricCounter{}) {
			c.JSON(http.StatusNotFound, nil)
		}
		floatValue, ok := metric.Get().(int64) 
		if !ok {
			c.JSON(http.StatusNotFound, nil)
		}
		m.Delta = &floatValue
		c.JSON(http.StatusOK, m)
	default:
		c.JSON(http.StatusBadRequest, nil)
	}
}
