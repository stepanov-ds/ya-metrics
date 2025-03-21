package main

import (
	"context"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
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

	r := gin.Default()
	var st storage.Storage
	var p *pgxpool.Pool

	ctx := context.Background()

	if server.IsDB {
		st = storage.NewDBStorage(ctx, storage.NewDBPool(ctx, *server.DatabaseDSN))
		p = st.(*storage.DBStorage).Pool
		defer st.(*storage.DBStorage).Pool.Close()
	} else {
		if *server.Restore {
			st = server.RestoreStorage()
		} else {
			st = storage.NewMemStorage(&sync.Map{})
		}
		go server.StoreInFile(st.(*storage.MemStorage))
	}
	router.Route(r, st, p)
	logger.Log.Info("main", zap.String("working with DB", strconv.FormatBool(server.IsDB)))

	if err := r.Run(*server.EndpointServer); err != nil {
		panic(err)
	}

	select {}
}
