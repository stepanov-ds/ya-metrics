package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"go.uber.org/zap"
)

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
