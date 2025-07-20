package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/collector"
	"github.com/stepanov-ds/ya-metrics/internal/config/agent"
	"github.com/stepanov-ds/ya-metrics/internal/sender"
)

var (
	Version   string
	BuildDate string
	Commit    string
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	wg := &sync.WaitGroup{}
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", Version, BuildDate, Commit)
	agent.ConfigAgent()
	var headers http.Header = make(map[string][]string)
	headers.Add("Content-Type", "application/json")
	sender := sender.NewHTTPSender(time.Second*10, headers, "http://"+*agent.EndpointAgent, *agent.RateLimit, agent.ReadPublicKey(*agent.CryptoKey).PublicKey.(*rsa.PublicKey))

	collector1 := collector.NewCollector(&sync.Map{})
	wg.Add(1)
	go collector1.Collect(ctx, wg, time.Duration(*agent.PollInterval)*time.Second, collector1.CollectMetrics)
	wg.Add(1)
	go sender.SendAll(ctx, wg, time.Duration(*agent.ReportInterval)*time.Second, collector1, true)

	collector2 := collector.NewCollector(&sync.Map{})
	wg.Add(1)
	go collector2.Collect(ctx, wg, time.Duration(*agent.PollInterval)*time.Second, collector2.CollectNewMetrics)
	wg.Add(1)
	go sender.SendAll(ctx, wg, time.Duration(*agent.ReportInterval)*time.Second, collector2, true)

	wg.Wait()
}
