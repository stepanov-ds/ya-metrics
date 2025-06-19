// Package handlers implements HTTP handlers for the metrics server.
//
// It includes:
// - Metric update and retrieval handlers
// - Health check and ping endpoints
// - Root endpoint to list all metrics
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"go.uber.org/zap"
)

// Updates handles bulk metric updates via JSON POST request.
//
// Expects a JSON array of utils.Metrics objects in the request body.
// Processes each metric and stores it using the provided storage.
// Supports transaction handling when using DBStorage.
func Updates(c *gin.Context, st storage.Storage) {
	var m []utils.Metrics
	// m := make([]utils.Metrics, 0, 30)
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
