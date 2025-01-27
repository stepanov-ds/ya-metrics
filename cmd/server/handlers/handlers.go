package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"github.com/stepanov-ds/ya-metrics/cmd/server/metric"
)


func Update(w http.ResponseWriter, r *http.Request, storage *metric.MemStorage) {
	if r.Method != http.MethodPost {
	     w.WriteHeader(http.StatusMethodNotAllowed)
	     return
	}

	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(path) < 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch strings.ToLower(path[1]) {
	case "gauge":
		if path[2] == "" || path[3] == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		gauge, err := strconv.ParseFloat(path[3], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric := metric.Metric{
			Gauge:     gauge,
			IsCounter: false,
		}
		storage.Storage[path[2]] = metric
	case "counter":
		if path[2] == "" || path[3] == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		counter, err := strconv.ParseInt(path[3], 0, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric := metric.Metric{
			Counter:   counter,
			IsCounter: true,
		}
		storage.Storage[path[2]] = metric
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}