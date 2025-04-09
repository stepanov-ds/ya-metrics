package server

import (
	"flag"
	"os"
	"strconv"

	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"go.uber.org/zap"
)

var (
	EndpointServer = flag.String("a", "localhost:8080", "endpoint")
	StoreInterval  = flag.Int("i", 300, "store interval")
	FileStorePath  = flag.String("f", "filestore.out", "file store path")
	Restore        = flag.Bool("r", true, "restore")
	DatabaseDSN    = flag.String("d", "", "database_DSN")
	IsDB           = false
	Key            = flag.String("k", "", "key")
)

func ConfigServer() {
	flag.Parse()
	address, found := os.LookupEnv("ADDRESS")
	if found {
		EndpointServer = &address
	}
	si, found := os.LookupEnv("STORE_INTERVAL")
	if found {
		i, err := strconv.Atoi(si)
		if err == nil && i >= 0 {
			StoreInterval = &i
		}
	}
	fsp, found := os.LookupEnv("FILE_STORAGE_PATH")
	if found {
		FileStorePath = &fsp
	}
	r, found := os.LookupEnv("RESTORE")
	if found {
		b, err := strconv.ParseBool(r)
		if err == nil {
			Restore = &b
		}
	}
	dsn, found := os.LookupEnv("DATABASE_DSN")
	if found {
		DatabaseDSN = &dsn
	}

	if *DatabaseDSN != "" {
		IsDB = true
	}
	k, found := os.LookupEnv("KEY")
	if found {
		Key = &k
	}
	*Key = ""
	logger.Log.Info("ConfigServer", 
		zap.String("EndpointServer", *EndpointServer),
		zap.Int("StoreInterval", *StoreInterval),
		zap.String("FileStorePath", *FileStorePath),
		zap.Bool("Restore", *Restore),
		zap.String("DatabaseDSN", *DatabaseDSN),
		zap.Bool("IsDB", IsDB),
		zap.String("Key", *Key),
	)
}
