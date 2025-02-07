package main

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/config"
	"github.com/stepanov-ds/ya-metrics/internal/handlers/router"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)
//metricstest -test.v -test.run=^TestIteration3[AB]$ -binary-path=cmd/server/server

func main() {
	config.ConfigServer()
	st := storage.NewMemStorage(&sync.Map{})
	r := gin.Default()
	router.Route(r, st)

	if err := r.Run(*config.EndpointS); err != nil {
		panic(err)
	}

}
