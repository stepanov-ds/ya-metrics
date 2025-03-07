package server

import (
	"encoding/json"
	"os"
	"time"

	"sync"

	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

func RestoreStorage() *storage.MemStorage {
	content, err := os.ReadFile(*FileStorePath)
	if err != nil {
		println(err.Error())
		return storage.NewMemStorage(&sync.Map{})
	}
	var metrics map[string]utils.Metrics
	err = json.Unmarshal(content, &metrics)
	if err != nil {
		println(err.Error())
		return storage.NewMemStorage(&sync.Map{})
	}
	var syncMap sync.Map
	for k, v := range metrics {
		// println(reflect.ValueOf(v).String())
		syncMap.Store(k, v)
	}
	return storage.NewMemStorage(&syncMap)
}

func storeInFile(s *storage.MemStorage) {
	jsonData, err := json.Marshal(s.GetAllMetrics())
	if err != nil {
		println(err.Error())
	}
	err = os.WriteFile(*FileStorePath, jsonData, os.FileMode(os.O_RDWR)|os.FileMode(os.O_CREATE)|os.FileMode(os.O_TRUNC))
	if err != nil {
		println(err.Error())
	}
}
func StoreInFile(s *storage.MemStorage) {
	for {
		time.Sleep(time.Second * time.Duration(*StoreInterval))
		storeInFile(s)
	}
}
