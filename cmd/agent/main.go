package main

import (
	"fmt"
	"net/http"
	"sync"
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
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", Version, BuildDate, Commit)
	agent.ConfigAgent()
	var headers http.Header = make(map[string][]string)
	headers.Add("Content-Type", "application/json")
	sender := sender.NewHTTPSender(time.Second*10, headers, "http://"+*agent.EndpointAgent, *agent.RateLimit)

	collector1 := collector.NewCollector(&sync.Map{})
	collector1.Collect(time.Duration(*agent.PollInterval)*time.Second, collector1.CollectMetrics)
	sender.Send(time.Duration(*agent.ReportInterval)*time.Second, collector1, true)

	collector2 := collector.NewCollector(&sync.Map{})
	collector2.Collect(time.Duration(*agent.PollInterval)*time.Second, collector2.CollectNewMetrics)
	sender.Send(time.Duration(*agent.ReportInterval)*time.Second, collector2, true)

	select {}

}
