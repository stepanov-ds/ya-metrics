package main

import (
	"net/http"
	"github.com/stepanov-ds/ya-metrics/cmd/server/handlers"
	"github.com/stepanov-ds/ya-metrics/cmd/server/storage"
)
//metricstest -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	st := storage.NewMemStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
        handlers.Update(w, r, st)
    })
	return http.ListenAndServe(`:8080`, mux)
}
