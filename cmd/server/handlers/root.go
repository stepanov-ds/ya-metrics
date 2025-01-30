package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/cmd/server/storage"
)

func Root(c *gin.Context, st storage.Storage) {
	jsonData, err := json.Marshal(st.GetAllMetrics())
	if err != nil {
		c.String(http.StatusOK, err.Error())
	}
	c.JSON(http.StatusOK, string(jsonData))
}