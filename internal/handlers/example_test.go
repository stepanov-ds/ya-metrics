package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/middlewares"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	r := gin.New()

	st := storage.NewMemStorage(&sync.Map{})
	r.Use(middlewares.Gzip())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.RedirectTrailingSlash = true

	r.Any("/update/:metric_type/:metric_name/:value", func(ctx *gin.Context) {
		Update(ctx, st)
	})
	r.Any("/update/:metric_type/:metric_name/:value/", func(ctx *gin.Context) {
		Update(ctx, st)
	})
	r.POST("/update", func(ctx *gin.Context) {
		Update(ctx, st)
	})
	r.GET("/value/:metric_type/:metric_name", func(ctx *gin.Context) {
		Value(ctx, st)
	})
	r.GET("/value/:metric_type/:metric_name/", func(ctx *gin.Context) {
		Value(ctx, st)
	})
	r.POST("/value", func(ctx *gin.Context) {
		Value(ctx, st)
	})
	r.GET("/", func(ctx *gin.Context) {
		Root(ctx, st)
	})
	r.GET("/ping", func(ctx *gin.Context) {
		Ping(ctx, nil)
	})
	r.POST("/updates", middlewares.HashCheck(), func(ctx *gin.Context) {
		Updates(ctx, st)
	})
	go func() {
		if err := r.Run("localhost:8080"); err != nil {
			panic(err)
		}
	}()
	time.Sleep(5 * time.Second)
}

// ExampleUpdate demonstrates how to update a gauge metric via JSON.
func ExampleUpdate() {
	client := &http.Client{}
	body :=
		`
	{
    	"id": "RandomValue",
    	"type": "gauge",
    	"value": 124.2
	}
	`
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/update/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		println(err.Error())
	}
	resp.Body.Close()
	fmt.Println(resp.Status)

	// Output:
	// 200 OK
}

// ExampleRoot demonstrates how to get all metrics
func ExampleRoot() {
	client := &http.Client{}
	body :=
		`
	{
    	"id": "RandomValue",
    	"type": "gauge",
    	"value": 124.2
	}
	`
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/update/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		println(err.Error())
	}
	resp.Body.Close()

	req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	resp, err = client.Do(req)
	if err != nil {
		println(err.Error())
	}
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(string(respBody))

	// Output:
	// {"RandomValue":{"value":124.2,"id":"RandomValue","type":"gauge"}}

}

// ExampleValue demonstrates how to get distinct metric
func ExampleValue() {
	client := &http.Client{}
	body :=
		`
	{
    	"id": "RandomValue",
    	"type": "gauge",
    	"value": 124.2
	}
	`
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/update/", strings.NewReader(body))
	// println(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		println(err.Error())
	}
	resp.Body.Close()

	req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/value/gauge/RandomValue", nil)
	resp, err = client.Do(req)
	if err != nil {
		println(err.Error())
	}
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(string(respBody))

	// Output:
	// 124.2

}
