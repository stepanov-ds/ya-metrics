package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/handlers"
	"github.com/stepanov-ds/ya-metrics/internal/middlewares"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)

func Route(r *gin.Engine, st *storage.MemStorage) {
	r.Use(middlewares.WithLogging())

	r.Any("/update/:metric_type/:metric_name/:value", func(ctx *gin.Context) {
		handlers.UpdateWithPath(ctx, st)
	})
	r.Any("/update/:metric_type/:metric_name/:value/", func(ctx *gin.Context) {
		handlers.UpdateWithPath(ctx, st)
	})
	r.POST("/update", func(ctx *gin.Context) {
		handlers.Update(ctx, st)
	})
	r.POST("/update/", func(ctx *gin.Context) {
		handlers.Update(ctx, st)
	})
	r.GET("/value/:metric_type/:metric_name", func(ctx *gin.Context) {
		handlers.ValueWithPath(ctx, st)
	})
	r.GET("/value/:metric_type/:metric_name/", func(ctx *gin.Context) {
		handlers.ValueWithPath(ctx, st)
	})
	r.POST("/value", func(ctx *gin.Context) {
		handlers.Value(ctx, st)
	})
	r.POST("/value/", func(ctx *gin.Context) {
		handlers.Value(ctx, st)
	})
	r.GET("/", func(ctx *gin.Context) {
		handlers.Root(ctx, st)
	})
}
