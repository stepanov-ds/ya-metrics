// Package handlers implements HTTP handlers for the metrics server.
//
// It includes:
// - Metric update and retrieval handlers
// - Health check and ping endpoints
// - Root endpoint to list all metrics
package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"go.uber.org/zap"
)

// Ping handles the /ping endpoint and checks database connectivity.
//
// If pool is nil or connection fails, responds with 500 Internal Server Error.
// On success, responds with 200 OK.
func Ping(c *gin.Context, pool *pgxpool.Pool) {
	if pool == nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	err := pool.Ping(context.Background())
	if err != nil {
		logger.Log.Info("Connection error: ",
			zap.String("Error", err.Error()),
		)
		c.String(http.StatusInternalServerError, "")
	} else {
		c.String(http.StatusOK, "")
	}
}
