package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/collector"
	"github.com/stepanov-ds/ya-metrics/internal/config/agent"
	"github.com/stepanov-ds/ya-metrics/internal/sender"
)

//metricstest --test.v --test.run=^TestIteration2[AB]*$ --source-path=. --agent-binary-path=cmd/agent/agent

func main() {
	agent.ConfigAgent()
	var headers http.Header = make(map[string][]string)
	collector := collector.NewCollector(&sync.Map{})
	headers.Add("Content-Type", "application/json")
	sender := sender.NewHTTPSender(time.Second*10, headers, "http://"+*agent.EndpointAgent)

	collector.Collect(time.Duration(*agent.PollInterval) * time.Second)
	sender.Send(time.Duration(*agent.ReportInterval)*time.Second, collector, true)

	select {}

}
