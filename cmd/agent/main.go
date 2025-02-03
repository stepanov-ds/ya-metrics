package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/collector"
	"github.com/stepanov-ds/ya-metrics/internal/sender"
)

var (
	endpoint       = flag.String("a", "localhost:8080", "endpoint")
	reportInterval = flag.Int("r", 10, "report interaval")
	pollInterval   = flag.Int("p", 2, "poll interval")
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
		for k, v := range collector.GetAllMetrics() {
			resp, err := sender.SendMetric(k, v)
			if err != nil {
				if resp != nil {
					print(resp.Body, err)
				}
			}
		}
		time.Sleep(interval)
	}
}

func main() {
	flag.Parse()
	ADDRESS, found := os.LookupEnv("ADDRESS")
	if found {
		endpoint = &ADDRESS
	}
	REPORT_INTERVAL, found := os.LookupEnv("REPORT_INTERVAL")
	if found {
		i, err := strconv.Atoi(REPORT_INTERVAL)
		if err != nil {
			reportInterval = &i
		}
	}
	POLL_INTERVAL, found := os.LookupEnv("POLL_INTERVAL")
	if found {
		i, err := strconv.Atoi(POLL_INTERVAL)
		if err != nil {
			pollInterval = &i
		}
	}
	var headers http.Header = make(map[string][]string)
	collector := collector.NewCollector()
	headers.Add("Content-Type", "text/plain")
	sender := sender.NewHttpSender(time.Second*10, headers, "http://"+*endpoint)

	go Collect(time.Duration(*pollInterval)*time.Second, collector)
	go Send(time.Duration(*reportInterval)*time.Second, collector, sender)

	select {}

}
