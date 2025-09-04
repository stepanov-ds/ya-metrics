package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/config/server"
	"github.com/stepanov-ds/ya-metrics/internal/grpcp"
	"github.com/stepanov-ds/ya-metrics/internal/handlers/router"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/stepanov-ds/ya-metrics/internal/grpcp/grpc_generated"
)

func main() {
	logger.Initialize("info")
	err := server.ConfigServer()
	if err != nil {
		logger.Log.Fatal("main", zap.Error(err))
	}

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
	logger.Log.Info("main", zap.String("working with DB", strconv.FormatBool(server.IsDB)))

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router.Route(r, st, p, server.ReadPrivateKey(*server.CryptoKey))

	

	if *server.Grpc {
		logger.Log.Info("main", zap.String("protocol", "gRPC"))
		lis, err := net.Listen("tcp", *server.EndpointServer)
		if err != nil {
			logger.Log.Fatal("failed to listen", zap.Error(err))
		}

		s := grpc.NewServer()
		pb.RegisterMetricsTunnelServer(s, &grpcp.GrpcHandler{Router: r})

		go func() {
			logger.Log.Info("gRPC server is running", zap.String("address", *server.EndpointServer))
			if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
				logger.Log.Fatal("failed to serve gRPC", zap.Error(err))
			}
		}()

		<-ctx.Done()
		logger.Log.Info("main", zap.String("shutdown", "shutting down gRPC server"))
		s.GracefulStop()

	} else {
		logger.Log.Info("main", zap.String("protocol", "http"))
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
}
