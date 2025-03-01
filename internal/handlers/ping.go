package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Ping(c *gin.Context, pool *pgxpool.Pool) {
	err := pool.Ping(context.Background())
	if err != nil {
		println(err.Error())
		c.String(http.StatusInternalServerError, "")
	} else {
		c.String(http.StatusOK, "")
	}
}