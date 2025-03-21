package agent

import (
	"flag"
	"os"
	"strconv"
)

var (
	EndpointAgent  = flag.String("a", "localhost:8080", "endpoint")
	ReportInterval = flag.Int("r", 10, "report interaval")
	PollInterval   = flag.Int("p", 2, "poll interval")
)

func ConfigAgent() {
	flag.Parse()
	address, found := os.LookupEnv("ADDRESS")
	if found {
		EndpointAgent = &address
	}
	ri, found := os.LookupEnv("REPORT_INTERVAL")
	if found {
		i, err := strconv.Atoi(ri)
		if err == nil && i >= 0 {
			ReportInterval = &i
		}
	}
	pi, found := os.LookupEnv("POLL_INTERVAL")
	if found {
		i, err := strconv.Atoi(pi)
		if err == nil && i >= 0 {
			PollInterval = &i
		}
	}
}
