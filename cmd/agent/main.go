package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/collector"
	"github.com/stepanov-ds/ya-metrics/internal/config/configagent"
	"github.com/stepanov-ds/ya-metrics/internal/sender"
)

//metricstest --test.v --test.run=^TestIteration2[AB]*$ --source-path=. --agent-binary-path=cmd/agent/agent

func main() {
	configagent.ConfigAgent()
	var headers http.Header = make(map[string][]string)
	collector := collector.NewCollector(&sync.Map{})
	headers.Add("Content-Type", "application/json")
	sender := sender.NewHTTPSender(time.Second*10, headers, "http://"+*configagent.Endpoint)

	collector.Collect(time.Duration(*configagent.PollInterval) * time.Second)
	sender.Send(time.Duration(*configagent.ReportInterval)*time.Second, collector)

	select {}

}
