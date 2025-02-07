package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/collector"
	"github.com/stepanov-ds/ya-metrics/internal/config"
	"github.com/stepanov-ds/ya-metrics/internal/sender"
)

//metricstest --test.v --test.run=^TestIteration2[AB]*$ --source-path=. --agent-binary-path=cmd/agent/agent

func main() {
	config.ConfigAgent()
	var headers http.Header = make(map[string][]string)
	collector := collector.NewCollector(&sync.Map{})
	headers.Add("Content-Type", "text/plain")
	sender := sender.NewHTTPSender(time.Second*10, headers, "http://"+*config.EndpointA)

	collector.Collect(time.Duration(*config.PollInterval) * time.Second)
	sender.Send(time.Duration(*config.ReportInterval)*time.Second, collector)

	select {}

}
