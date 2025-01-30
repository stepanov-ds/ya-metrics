package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/cmd/server/storage"
	"github.com/stepanov-ds/ya-metrics/pkg/utils"
)

func Update(c *gin.Context, st storage.Storage) {
	metricType := c.Param("metric_type")
	metricName := c.Param("metric_name")
	metricValue := c.Param("value")

	if metricType == "" || metricName == "" || metricValue == "" {
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
		metric := utils.Metric{
			Gauge:     gauge,
			IsCounter: false,
		}
		st.SetMetric(metricName, metric)
	case "counter":
		counter, err := strconv.ParseInt(metricValue, 0, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		oldMetricValue, found := st.GetMetric(metricName)
		metric := utils.Metric{}
		if found {
			if oldMetricValue.IsCounter {
				metric = utils.Metric{
					Counter:   counter + oldMetricValue.Counter,
					IsCounter: true,
				}
			} else {
				metric = utils.Metric{
					Counter:   counter,
					IsCounter: true,
				}
			}
		} else {
			metric = utils.Metric{
				Counter:   counter,
				IsCounter: true,
			}
		}
		st.SetMetric(metricName, metric)
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Data(http.StatusOK, "", nil)
}
