// Package server implements configuration logic for the metrics server.
//
// The server is responsible for:
// - Parsing command-line flags and environment variables
// - Setting up storage (file or database)
// - Configuring network endpoint
// - Optional payload signature verification with a key
package server

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"go.uber.org/zap"
)

var (
	// EndpointServer holds the server address in the format "host:port".
	// Can be set via flag "-a" or env var "ADDRESS".
	EndpointServer = flag.String("a", "localhost:8080", "endpoint")
	// StoreInterval defines how often (in seconds) the server saves metrics to file.
	// A value of 0 means synchronous writes.
	// Can be set via flag "-i" or env var "STORE_INTERVAL".
	StoreInterval = flag.Int("i", 300, "store interval")
	// FileStorePath specifies the path to the file where metrics are stored.
	// Can be set via flag "-f" or env var "FILE_STORAGE_PATH".
	FileStorePath = flag.String("f", "filestore.out", "file store path")
	// Restore defines whether metrics should be restored from the file on startup.
	// Can be set via flag "-r" or env var "RESTORE".
	Restore = flag.Bool("r", true, "restore")
	// DatabaseDSN contains the DSN (Data Source Name) for connecting to a database.
	// If empty, file-based storage is used.
	// Can be set via flag "-d" or env var "DATABASE_DSN".
	DatabaseDSN = flag.String("d", "", "database_DSN")
	// IsDB indicates whether the server is configured to use a database.
	// This is derived from the presence of DatabaseDSN.
	IsDB = false
	// Key holds an optional signing key used to verify metric payloads.
	// Can be set via flag "-k" or env var "KEY".
	Key = flag.String("k", "", "key")
	// Crypto hold a private key
	CryptoKey  = flag.String("y", "private_key.pem", "crypto key")
	ConfigFile = flag.String("c", "", "config file")
	Loaded     = false
)

// ConfigServer parses command-line flags and environment variables
// to configure the server at runtime.
//
// It respects the following precedence:
// 1. Command-line flags override environment variables.
// 2. Environment variables override defaults.
//
// After parsing, it logs the final configuration using zap.Logger.
func ConfigServer() error {
	var CryptoKeyVar string
	var ConfigFileVar string
	flag.StringVar(&CryptoKeyVar, "crypto-key", "cert.pem", "crypto key")
	flag.StringVar(&ConfigFileVar, "config", "", "config")
	flag.Parse()

	if CryptoKeyVar != "" {
		*CryptoKey = CryptoKeyVar
	}
	if ConfigFileVar != "" {
		*ConfigFile = ConfigFileVar
	}

	if *ConfigFile == "" {
		*ConfigFile = os.Getenv("CONFIG")
	}

	err := LoadConfigFile()
	if err != nil {
		return err
	}

	// if Loaded {
	// 	*EndpointServer, *FileStorePath , *DatabaseDSN, *CryptoKey, *StoreInterval, IsDB = a, b, c, d, e, f
	// }

	loadFromEnv()

	logger.Log.Info("ConfigServer",
		zap.String("EndpointServer", *EndpointServer),
		zap.Int("StoreInterval", *StoreInterval),
		zap.String("FileStorePath", *FileStorePath),
		zap.Bool("Restore", *Restore),
		zap.String("DatabaseDSN", *DatabaseDSN),
		zap.Bool("IsDB", IsDB),
		zap.String("Key", *Key),
	)
	return nil
}

type Config struct {
	Address       string `json:"address,omitempty"`
	StoreFile     string `json:"store_file,omitempty"`
	DatabaseDSN   string `json:"database_dsn,omitempty"`
	CryptoKey     string `json:"crypto_key,omitempty"`
	StoreInterval string `json:"store_interval,omitempty"`
	Restore       bool   `json:"restore,omitempty"`
}

func loadFromEnv() {
	// Override with environment variables if present
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
	cr, found := os.LookupEnv("CRYPTO_KEY")
	if found {
		CryptoKey = &cr
	}
}

func LoadConfigFile() error {
	cfg := &Config{}
	a, b, c, d, e, f, ab, bb, cb, db, eb, fb := storeParsed()

	if *ConfigFile != "" {
		file, err := os.Open(*ConfigFile)
		if err != nil {
			return err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err = decoder.Decode(cfg); err != nil {
			return err
		}
		*EndpointServer = cfg.Address
		dur, err := time.ParseDuration(cfg.StoreInterval)
		if err != nil {
			return err
		}
		*StoreInterval = int(dur.Seconds())
		*FileStorePath = cfg.StoreFile
		*Restore = cfg.Restore
		*DatabaseDSN = cfg.DatabaseDSN
		if cfg.DatabaseDSN != "" {
			IsDB = true
		}
		*CryptoKey = cfg.CryptoKey
		Loaded = true
	}
	checkLoaded(a, b, c, d, e, f, ab, bb, cb, db, eb, fb)
	return nil
}
func storeParsed() (string, int, string, string, string, bool, bool, bool, bool, bool, bool, bool) {
	var a, c, d, e string
	var b int
	var f bool

	ab, bb, cb, db, eb, fb := false, false, false, false, false, false
	if *EndpointServer != "localhost:8080" {
		a = *EndpointServer
		ab = true
	}
	if *StoreInterval != 300 {
		b = *StoreInterval
		bb = true
	}
	if *FileStorePath != "filestore.out" {
		c = *FileStorePath
		cb = true
	}
	if *CryptoKey != "private_key.pem" {
		d = *CryptoKey
		db = true
	}
	if *DatabaseDSN != "" {
		e = *DatabaseDSN
		eb = true
	}
	if IsDB {
		f = IsDB
		fb = true
	}
	return a, b, c, d, e, f, ab, bb, cb, db, eb, fb
}

func checkLoaded(a string, b int, c, d, e string, f, ab, bb, cb, db, eb, fb bool) {
	if Loaded {
		if ab {
			*EndpointServer = a
		}
		if bb {
			*StoreInterval = b
		}
		if cb {
			*FileStorePath = c
		}
		if db {
			*CryptoKey = d
		}
		if eb {
			*DatabaseDSN = e
		}
		if fb {
			IsDB = f
		}
		// *EndpointAgent, *ReportInterval , *PollInterval, *CryptoKey = a, b, c, d
	}
}
