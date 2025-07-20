package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/config/server"
	"github.com/stepanov-ds/ya-metrics/internal/handlers/router"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"go.uber.org/zap"
)

func main() {
	logger.Initialize("info")
	server.ConfigServer()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	var st storage.Storage
	var p *pgxpool.Pool

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	if server.IsDB {
		st = storage.NewDBStorage(ctx, storage.NewDBPool(ctx, *server.DatabaseDSN))
		p = st.(*storage.DBStorage).Pool.(*pgxpool.Pool)
		defer st.(*storage.DBStorage).Pool.(*pgxpool.Pool).Close()
	} else {
		if *server.Restore {
			st = server.RestoreStorage()
		} else {
			st = storage.NewMemStorage(&sync.Map{})
		}
		go server.StoreInFile(st.(*storage.MemStorage))
	}
	router.Route(r, st, p, server.ReadPrivateKey(*server.CryptoKey))
	logger.Log.Info("main", zap.String("working with DB", strconv.FormatBool(server.IsDB)))

	srv := &http.Server{
		Addr:    *server.EndpointServer,
		Handler: r.Handler(),
	}

	quit := make(chan os.Signal, 1)
	idleConnsClosed := make(chan struct{})
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Log.Info("main", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Info("main", zap.Error(err))
	}
	<-idleConnsClosed
}
