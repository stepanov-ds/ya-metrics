package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/config/server"
	"github.com/stepanov-ds/ya-metrics/internal/handlers"
	"github.com/stepanov-ds/ya-metrics/internal/middlewares"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)

func Route(r *gin.Engine, st storage.Storage, pool *pgxpool.Pool) {
	r.Use(middlewares.Gzip())
	r.Use(middlewares.WithLogging())
	if *server.Key != "" {
		r.Use(middlewares.HashCheck())
	}

	r.RedirectTrailingSlash = true

	r.Any("/update/:metric_type/:metric_name/:value", func(ctx *gin.Context) {
		handlers.Update(ctx, st)
	})
	r.Any("/update/:metric_type/:metric_name/:value/", func(ctx *gin.Context) {
		handlers.Update(ctx, st)
	})
	r.POST("/update", func(ctx *gin.Context) {
		handlers.Update(ctx, st)
	})
	// r.POST("/update/", func(ctx *gin.Context) {
	// 	handlers.Update(ctx, st)
	// })
	r.GET("/value/:metric_type/:metric_name", func(ctx *gin.Context) {
		handlers.Value(ctx, st)
	})
	r.GET("/value/:metric_type/:metric_name/", func(ctx *gin.Context) {
		handlers.Value(ctx, st)
	})
	r.POST("/value", func(ctx *gin.Context) {
		handlers.Value(ctx, st)
	})
	// r.POST("/value/", func(ctx *gin.Context) {
	// 	handlers.Value(ctx, st)
	// })
	r.GET("/", func(ctx *gin.Context) {
		handlers.Root(ctx, st)
	})
	r.GET("/ping", func(ctx *gin.Context) {
		handlers.Ping(ctx, pool)
	})
	// r.GET("/ping/", func(ctx *gin.Context) {
	// 	handlers.Ping(ctx, pool)
	// })
	r.POST("/updates", func(ctx *gin.Context) {
		handlers.Updates(ctx, st)
	})
	// r.POST("/updates/", func(ctx *gin.Context) {
	// 	handlers.Updates(ctx, st)
	// })
}
