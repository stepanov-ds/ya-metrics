package main

import (
	"context"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/config/server"
	"github.com/stepanov-ds/ya-metrics/internal/handlers/router"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"go.uber.org/zap"
)

//metricstest -test.v -test.run=^TestIteration3[AB]$ -binary-path=cmd/server/server

func main() {
	logger.Initialize("info")
	server.ConfigServer()
	var st storage.Storage

	if server.IsDb {
		st = storage.NewDbStorage(storage.NewDbPool(context.Background(), *server.Database_DSN))
		defer st.(*storage.DbStorage).Pool.Close()
	} else {
		if *server.Restore {
			st = server.RestoreStorage()
		} else {
			st = storage.NewMemStorage(&sync.Map{})
		}
		go server.StoreInFile(st.(*storage.MemStorage))
	}
	logger.Log.Info("main", zap.String("working with DB", strconv.FormatBool(server.IsDb)))

	r := gin.Default()
	router.MainRoute(r, st)

	if err := r.Run(*server.EndpointServer); err != nil {
		panic(err)
	}

	select {}
}
