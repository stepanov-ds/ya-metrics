package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/cmd/server/handlers"
	"github.com/stepanov-ds/ya-metrics/cmd/server/storage"
)

//metricstest -test.v -test.run=^TestIteration3[AB]$ -binary-path=cmd/server/server

func main() {
	st := storage.NewMemStorage()
	r := gin.Default()

	r.Any("/update/:metric_type/:metric_name/:value", func(ctx *gin.Context) {
		st.LockMutex()
		defer st.UnlockMutex()
		handlers.Update(ctx, st)
	})
	r.Any("/update/:metric_type/:metric_name/:value/", func(ctx *gin.Context) {
		st.LockMutex()
		defer st.UnlockMutex()
		handlers.Update(ctx, st)
	})
	r.GET("/value/:metric_type/:metric_name", func(ctx *gin.Context) {
		st.LockMutex()
		defer st.UnlockMutex()
		handlers.Value(ctx, st)
	})
	r.GET("/value/:metric_type/:metric_name/", func(ctx *gin.Context) {
		st.LockMutex()
		defer st.UnlockMutex()
		handlers.Value(ctx, st)
	})
	r.GET("/", func(ctx *gin.Context) {
		st.LockMutex()
		defer st.UnlockMutex()
		handlers.Root(ctx, st)
	})
	

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

}
