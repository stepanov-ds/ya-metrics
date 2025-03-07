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
	var m *utils.Metrics

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
				metricValue = strconv.FormatInt(*m.Delta, 10)
			} else if strings.ToLower(metricType) == "gauge" {
				metricValue = strconv.FormatFloat(*m.Value, 'f', -1, 64)
			} else {
				c.AbortWithStatus(http.StatusBadRequest)
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
		st.SetMetric(metricName, v, strings.ToLower(metricType) == "counter")
	} else if strings.ToLower(metricType) == "counter" {
		v, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		st.SetMetric(metricName, v, strings.ToLower(metricType) == "counter")
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Data(http.StatusOK, "", nil)
}

func UpdateWithJSON(c *gin.Context, st storage.Storage) *utils.Metrics { //где-то тут хуета
	var m utils.Metrics
	if err := c.ShouldBindJSON(&m); err != nil {
		return &utils.Metrics{}
	}
	return &m

}
