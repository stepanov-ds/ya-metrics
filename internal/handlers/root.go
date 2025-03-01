package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)

func Root(c *gin.Context, st storage.Storage) {
	jsonData, err := json.Marshal(st.GetAllMetrics())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Writer.Header().Add("Content-Type","text/html")
	c.String(http.StatusOK, string(jsonData))
}
