package main

import (
	"net/http"
	"strconv"
	"strings"
)

var storage MemStorage = MemStorage {
	Storage: make(map[string]Metric),
}

type Metric struct {
	Counter   int64
	Gauge     float64
	IsCounter bool
}

type MemStorage struct {
	Storage map[string]Metric
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", update)
	return http.ListenAndServe(`:8080`, mux)
}

func update(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	//     w.WriteHeader(http.StatusMethodNotAllowed)
	//     return
	// }

	path := strings.Split(strings.Trim(strings.ToLower(r.URL.Path), "/"), "/")

	//checking metric type
	if len(path) >= 2 {
		switch path[1] {
		case "gauge":
			//checking metric name
			if len(path) >= 3 {
				if path[2] == "" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				//checking metric value
				if len(path) >=4 {
					if path[3] == "" {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					gauge, err := strconv.ParseFloat(path[3], 64)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					metric := Metric{
						Gauge: gauge,
						IsCounter: false,
					}
					storage.Storage[path[2]] = metric
				} else {
					w.WriteHeader(http.StatusNotFound)
						return
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		case "counter":
			//checking metric name
			if len(path) >= 3 {
				if path[2] == "" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				//checking metric value
				if len(path) >=4 {
					if path[3] == "" {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					counter, err := strconv.ParseInt(path[3], 0, 64)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					metric := Metric{
						Counter: counter,
						IsCounter: true,
					}
					storage.Storage[path[2]] = metric
				} else {
					w.WriteHeader(http.StatusNotFound)
						return
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)

	// w.Write()

}
