package main

import (
	"net/http"
	"time"

	"github.com/stepanov-ds/ya-metrics/cmd/agent/collector"
	"github.com/stepanov-ds/ya-metrics/cmd/agent/sender"
)

//metricstest --test.v --test.run=^TestIteration2[AB]*$ --source-path=. --agent-binary-path=cmd/agent/agent

func Collect(interval time.Duration, collector *collector.Collector) {
	for {
		collector.CollectMetrics()
		time.Sleep(interval)
	}
}

func Send(interval time.Duration, collector *collector.Collector, sender sender.Sender) {
	for {
		for k, v := range collector.Metrics {
			sender.SendMetric(k, v)
		}
		time.Sleep(interval)
	}
}

func main() {
	baseUrl := "http://localhost:8080"
	var headers http.Header = make(map[string][]string)
	headers.Add("Content-Type", "text/plain")
	collector := collector.NewCollector()
	sender := sender.NewHttpSender(time.Second*10, headers, baseUrl)

	pollInterval := time.Second * 2
	reportInterval := time.Second * 10

	go Collect(pollInterval, collector)
	go Send(reportInterval, collector, &sender)

	select {}

}
