// Package agent implements configuration logic for the metrics agent.
//
// The agent is responsible for:
// - Parsing command-line flags
// - Reading environment variables
// - Setting up connection endpoint and reporting intervals
package agent

import (
	"flag"
	"os"
	"strconv"
)

var (
	// EndpointAgent holds the server address in the format "host:port".
    // Can be set via flag "-a" or env var "ADDRESS".
	EndpointAgent  = flag.String("a", "localhost:8080", "endpoint")
	// ReportInterval defines how often (in seconds) the agent sends metrics to the server.
    // Can be set via flag "-r" or env var "REPORT_INTERVAL".
	ReportInterval = flag.Int("r", 10, "report interaval")
	// PollInterval defines how often (in seconds) the agent collects new metrics locally.
    // Can be set via flag "-p" or env var "POLL_INTERVAL".
	PollInterval   = flag.Int("p", 2, "poll interval")
	// Key holds an optional signing key used to calculate hash of the metric payload.
    // Can be set via flag "-k" or env var "KEY".
	Key            = flag.String("k", "", "key")
	// RateLimit defines maximum number of concurrent requests to the server.
    // Can be set via flag "-l" or env var "RATE_LIMIT".
	RateLimit      = flag.Int("l", 1, "rate limit")
)

// ConfigAgent parses command-line flags and environment variables
// to configure agent behavior at runtime.
//
// It respects the following precedence:
// 1. Command-line flags override environment variables.
// 2. Environment variables override defaults.
func ConfigAgent() {
	flag.Parse()

	// Override with environment variables if present
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
	
	// Print current config values
	println("EndpointAgent=", *EndpointAgent)
	println("ReportInterval=", *ReportInterval)
	println("PollInterval=", *PollInterval)
	println("Key=", *Key)
	println("RateLimit=", *RateLimit)
}
