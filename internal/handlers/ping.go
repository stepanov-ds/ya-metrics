package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"go.uber.org/zap"
)

func Ping(c *gin.Context, pool *pgxpool.Pool) {
	err := pool.Ping(context.Background())
	if err != nil {
		println(err.Error())
		logger.Log.Info("Connection configuration: ", 
			zap.String("Host", pool.Config().ConnConfig.Host),
			zap.String("Port", strconv.FormatInt(int64(pool.Config().ConnConfig.Port), 10)),
			zap.String("User", pool.Config().ConnConfig.User),
			zap.String("Database", pool.Config().ConnConfig.Database),
		)
		c.String(http.StatusInternalServerError, "")
	} else {
		c.String(http.StatusOK, "")
	}
}