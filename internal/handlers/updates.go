package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"go.uber.org/zap"
)

func Updates(c *gin.Context, st storage.Storage) {
	var m []utils.Metrics
	if err := c.ShouldBindBodyWithJSON(&m); err != nil {
		logger.Log.Error("Updates", zap.String("error while unmarshal body", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := c.Request.Context()

	var isDB bool
	if _, isDB := st.(*storage.DBStorage); isDB {
		st.(*storage.DBStorage).BeginTransaction(ctx)
		defer st.(*storage.DBStorage).RollbackTransaction(ctx)
	}

	for _, item := range m {
		if item.MType == "counter" {
			st.SetMetric(ctx, item.ID, item.Delta, true)
		} else if item.MType == "gauge" {
			st.SetMetric(ctx, item.ID, item.Value, false)
		}
	}

	if isDB {
		st.(*storage.DBStorage).CommitTransaction(ctx)
	}

	c.Data(http.StatusOK, "", nil)
}
