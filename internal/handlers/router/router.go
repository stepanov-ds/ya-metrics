// Package router sets up HTTP routes and middleware for the metrics server.
//
// It defines:
// - Metric update and query endpoints
// - Middleware stack (logging, compression, hash validation)
// - Profiling routes via pprof
package router

import (
	// "net/http"

	"crypto/rsa"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/handlers"
	"github.com/stepanov-ds/ya-metrics/internal/middlewares"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)

// Route registers all HTTP handlers and middleware for the Gin engine.
//
// Registers:
// - Gzip compression middleware
// - Request logging middleware
// - Hash validation middleware (optional)
// - Metric update and value retrieval endpoints
// - Pprof profiling routes
func Route(r *gin.Engine, st storage.Storage, pool *pgxpool.Pool, privateKey *rsa.PrivateKey) {
	r.Use(middlewares.Gzip())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(middlewares.Crypto(privateKey))

	r.Use(middlewares.WithLogging())

	// if *server.Key != "" {
	// 	r.Use(middlewares.HashCheck())
	// }

	r.RedirectTrailingSlash = true

	// Update metric by URL path
	r.Any("/update/:metric_type/:metric_name/:value", func(ctx *gin.Context) {
		handlers.Update(ctx, st)
	})
	r.Any("/update/:metric_type/:metric_name/:value/", func(ctx *gin.Context) {
		handlers.Update(ctx, st)
	})

	// Update metric via JSON body
	r.POST("/update", func(ctx *gin.Context) {
		handlers.Update(ctx, st)
	})
	// r.POST("/update/", func(ctx *gin.Context) {
	// 	handlers.Update(ctx, st)
	// })

	// Get metric value by name and type
	r.GET("/value/:metric_type/:metric_name", func(ctx *gin.Context) {
		handlers.Value(ctx, st)
	})
	r.GET("/value/:metric_type/:metric_name/", func(ctx *gin.Context) {
		handlers.Value(ctx, st)
	})

	// Get metric value via JSON body
	r.POST("/value", func(ctx *gin.Context) {
		handlers.Value(ctx, st)
	})
	// r.POST("/value/", func(ctx *gin.Context) {
	// 	handlers.Value(ctx, st)
	// })

	// Root endpoint to list all metrics
	r.GET("/", func(ctx *gin.Context) {
		handlers.Root(ctx, st)
	})

	// Database ping endpoint
	r.GET("/ping", func(ctx *gin.Context) {
		handlers.Ping(ctx, pool)
	})
	// r.GET("/ping/", func(ctx *gin.Context) {
	// 	handlers.Ping(ctx, pool)
	// })

	// Bulk updates with hash validation
	r.POST("/updates", middlewares.HashCheck(), func(ctx *gin.Context) {
		handlers.Updates(ctx, st)
	})
	// r.POST("/updates/", func(ctx *gin.Context) {
	// 	handlers.Updates(ctx, st)
	// })

	// Register pprof profiling routes under /debug/pprof/*
	pprof.Register(r)
}
