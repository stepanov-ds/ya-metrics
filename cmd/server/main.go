package main

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/config/server"
	"github.com/stepanov-ds/ya-metrics/internal/handlers/router"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)

//metricstest -test.v -test.run=^TestIteration3[AB]$ -binary-path=cmd/server/server

func main() {
	pool, err := pgxpool.New(context.Background(), *server.Database_DSN)
	if err != nil {
		println(err.Error())
	}
	defer pool.Close()
	logger.Initialize("info")
	server.ConfigServer()

	var st *storage.MemStorage
	if *server.Restore {
		st = server.RestoreStorage()
	} else {
		st = storage.NewMemStorage(&sync.Map{})
	}
	r := gin.Default()
	router.Route(r, st, pool)
	go server.StoreInFile(st)
	if err := r.Run(*server.EndpointServer); err != nil {
		panic(err)
	}

	select{}
}
