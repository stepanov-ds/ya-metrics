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
	Key            = flag.String("k", "", "key")
	RateLimit      = flag.Int("l", 1, "rate limit")
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
	k, found := os.LookupEnv("KEY")
	if found {
		Key = &k
	}
	rl, found := os.LookupEnv("RATE_LIMIT")
	if found {
		i, err := strconv.Atoi(rl)
		if err == nil && i >= 0 {
			RateLimit = &i
		}
	}
	println("EndpointAgent=", *EndpointAgent)
	println("ReportInterval=", *ReportInterval)
	println("PollInterval=", *PollInterval)
	println("Key=", *Key)
	println("RateLimit=", *RateLimit)
}
