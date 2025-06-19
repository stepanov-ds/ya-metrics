// Package server implements logic for restoring and persisting metrics storage.
//
// It supports:
// - Loading metrics from a file on startup
// - Periodically saving current metrics to a file
package server

import (
	"encoding/json"
	"os"
	"time"

	"sync"

	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"go.uber.org/zap"
)

// RestoreStorage loads metrics from the storage file (if exists) and returns an initialized MemStorage.
//
// If the file is missing or contains invalid data, it returns a new empty storage.
func RestoreStorage() *storage.MemStorage {
	content, err := os.ReadFile(*FileStorePath)
	if err != nil {
		logger.Log.Error("RestoreStorage", zap.String("error while reading file", err.Error()))
		return storage.NewMemStorage(&sync.Map{})
	}
	var metrics map[string]utils.Metrics
	err = json.Unmarshal(content, &metrics)
	if err != nil {
		logger.Log.Error("RestoreStorage", zap.String("error while unmarshal file", err.Error()))
		return storage.NewMemStorage(&sync.Map{})
	}
	var syncMap sync.Map
	for k, v := range metrics {
		// println(reflect.ValueOf(v).String())
		syncMap.Store(k, v)
	}
	return storage.NewMemStorage(&syncMap)
}

// storeInFile saves all current metrics from MemStorage to the storage file in JSON format.
//
// If marshaling or writing fails, logs an error using zap.Logger.
func storeInFile(s *storage.MemStorage) {
	jsonData, err := json.Marshal(s.GetAllMetrics())
	if err != nil {
		logger.Log.Error("storeInFile", zap.String("error while marshal metrics", err.Error()))
	}
	err = os.WriteFile(*FileStorePath, jsonData, os.FileMode(os.O_RDWR)|os.FileMode(os.O_CREATE)|os.FileMode(os.O_TRUNC))
	if err != nil {
		logger.Log.Error("storeInFile", zap.String("error while writing file", err.Error()))
	}
}

// StoreInFile starts a background loop that periodically saves metrics to disk.
//
// Interval is defined by StoreInterval (in seconds).
func StoreInFile(s *storage.MemStorage) {
	for {
		time.Sleep(time.Second * time.Duration(*StoreInterval))
		storeInFile(s)
	}
}
