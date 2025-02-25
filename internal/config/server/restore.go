package server

import (
	"encoding/json"
	"os"
	"reflect"
	"sync"

	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

func RestoreStorage() *storage.MemStorage {
	content, err := os.ReadFile(*FileStorePath)
	if err != nil {
		println("READ_ERR")
		return storage.NewMemStorage(&sync.Map{})
	}
	var metrics map[string]utils.Metrics
	err = json.Unmarshal(content, &metrics)
	if err != nil {
		println("UNMARSHAL_ERR")
		return storage.NewMemStorage(&sync.Map{})
	}
	var syncMap sync.Map
	for k, v := range metrics {
		println(reflect.ValueOf(v).String())
		syncMap.Store(k, v)
	}
	return storage.NewMemStorage(&syncMap)
}
