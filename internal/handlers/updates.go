package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

	if _, ok := st.(*storage.DBStorage); ok {
		ctx := c.Request.Context()
		tx, err := st.(*storage.DBStorage).Pool.Begin(ctx)
		if err != nil {
			logger.Log.Error("Updates", zap.String("error while starting transaction", err.Error()))
		}
		defer func() {
			if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
				logger.Log.Error("Updates", zap.String("error while rollback", err.Error()))
			}
		}()
	}

	for _, item := range m {
		if item.MType == "counter" {
			st.SetMetric(item.ID, item.Delta, true)
		} else if item.MType == "gauge" {
			st.SetMetric(item.ID, item.Value, false)
		}
	}

	c.Data(http.StatusOK, "", nil)
}
