package main

import (
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/handlers"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)

var (
	endpoint = flag.String("a", "localhost:8080", "endpoint")
)

//metricstest -test.v -test.run=^TestIteration3[AB]$ -binary-path=cmd/server/server

func main() {
	flag.Parse()
	ADDRESS, found := os.LookupEnv("ADDRESS")
	if found {
		endpoint = &ADDRESS
	}
	st := storage.NewMemStorage()
	r := gin.Default()

	r.Any("/update/:metric_type/:metric_name/:value", func(ctx *gin.Context) {
		handlers.Update(ctx, st)
	})
	r.Any("/update/:metric_type/:metric_name/:value/", func(ctx *gin.Context) {
		handlers.Update(ctx, st)
	})
	r.GET("/value/:metric_type/:metric_name", func(ctx *gin.Context) {
		handlers.Value(ctx, st)
	})
	r.GET("/value/:metric_type/:metric_name/", func(ctx *gin.Context) {
		handlers.Value(ctx, st)
	})
	r.GET("/", func(ctx *gin.Context) {
		handlers.Root(ctx, st)
	})

	if err := r.Run(*endpoint); err != nil {
		panic(err)
	}

}
