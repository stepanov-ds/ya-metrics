package main

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/config/server"
	"github.com/stepanov-ds/ya-metrics/internal/handlers/router"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)

//metricstest -test.v -test.run=^TestIteration3[AB]$ -binary-path=cmd/server/server

func main() {
	logger.Initialize("info")
	server.ConfigServer()

	var st *storage.MemStorage
	if *server.Restore {
		st = server.RestoreStorage()
	} else {
		st = storage.NewMemStorage(&sync.Map{})
	}
	r := gin.Default()
	router.Route(r, st)

	if err := r.Run(*server.EndpointServer); err != nil {
		panic(err)
	}
}
