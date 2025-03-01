package server

import (
	"flag"
	"os"
	"strconv"
	"time"
)

var (
	EndpointServer = flag.String("a", "localhost:8080", "endpoint")
	StoreInterval  = flag.Int("i", 300, "store interval")
	FileStorePath  = flag.String("f", "filestore.out", "file store path")
	Restore        = flag.Bool("r", true, "restore")
	LastFileWrite  = time.Now()
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
}
