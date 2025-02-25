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

	defer c.Request.Body.Close()
	if err := c.ShouldBindJSON(&m); err == nil {
		c.JSON(http.StatusBadRequest, nil)
		return m
	} else {
		println(err.Error())
		return utils.Metrics{}
	}
}
	// if m.ID == "" {
	// 	c.JSON(http.StatusBadRequest, nil)
	// 	return
	// }
	// metric, found := st.GetMetric(m.ID)
	// if found {
	// 	if strings.ToLower(m.ID) == strings.ToLower(metric.MType) {
	// 		c.String(http.StatusOK, fmt.Sprintf("%v", metric.Get()))
	// 	} else {
	// 		c.String(http.StatusNotFound, "")
	// 	}
	// } else {
	// 	c.String(http.StatusNotFound, "")
	// }

	// switch strings.ToLower(m.MType) {
	// case "gauge":
	// 	metric, found := st.GetMetric(m.ID)
	// 	if !found {
	// 		c.JSON(http.StatusNotFound, nil)
	// 		return
	// 	}
	// 	if reflect.TypeOf(metric) != reflect.TypeOf(&utils.MetricGauge{}) {
	// 		c.JSON(http.StatusNotFound, nil)
	// 		return
	// 	}
	// 	floatValue, ok := metric.Get().(float64)
	// 	if !ok {
	// 		c.JSON(http.StatusNotFound, nil)
	// 		return
	// 	}
	// 	m.Value = &floatValue
	// 	c.JSON(http.StatusOK, m)
	// 	return
	// case "counter":
	// 	metric, found := st.GetMetric(m.ID)
	// 	if !found {
	// 		c.JSON(http.StatusNotFound, nil)
	// 		return
	// 	}
	// 	if reflect.TypeOf(metric) != reflect.TypeOf(&utils.MetricCounter{}) {
	// 		c.JSON(http.StatusNotFound, nil)
	// 		return
	// 	}
	// 	floatValue, ok := metric.Get().(int64)
	// 	if !ok {
	// 		c.JSON(http.StatusNotFound, nil)
	// 		return
	// 	}
	// 	m.Delta = &floatValue
	// 	c.JSON(http.StatusOK, m)
	// 	return
	// default:
	// 	c.JSON(http.StatusBadRequest, nil)
	// 	return
	// }

