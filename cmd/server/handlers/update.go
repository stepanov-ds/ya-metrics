package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"github.com/stepanov-ds/ya-metrics/cmd/server/storage"
	"github.com/stepanov-ds/ya-metrics/pkg/utils"
)


func Update(w http.ResponseWriter, r *http.Request, repo storage.Repositories) {
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
		metric := utils.Metric{
			Gauge:     gauge,
			IsCounter: false,
		}
		repo.SetMetric(path[2], metric)
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
		metric := utils.Metric{
			Counter:   counter,
			IsCounter: true,
		}
		repo.SetMetric(path[2], metric)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}