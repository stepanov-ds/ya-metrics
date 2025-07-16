// Package agent implements configuration logic for the metrics agent.
//
// The agent is responsible for:
// - Parsing command-line flags
// - Reading environment variables
// - Setting up connection endpoint and reporting intervals
package agent

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"time"
)

var (
	// Crypto hold a public key
	CryptoKey  = flag.String("y", "cert.pem", "crypto key")
	ConfigPath = flag.String("c", "", "config")
	// EndpointAgent holds the server address in the format "host:port".
	// Can be set via flag "-a" or env var "ADDRESS".
	EndpointAgent = flag.String("a", "localhost:8080", "endpoint")
	// ReportInterval defines how often (in seconds) the agent sends metrics to the server.
	// Can be set via flag "-r" or env var "REPORT_INTERVAL".
	ReportInterval = flag.Int("r", 10, "report interaval")
	// PollInterval defines how often (in seconds) the agent collects new metrics locally.
	// Can be set via flag "-p" or env var "POLL_INTERVAL".
	PollInterval = flag.Int("p", 2, "poll interval")
	// Key holds an optional signing key used to calculate hash of the metric payload.
	// Can be set via flag "-k" or env var "KEY".
	Key = flag.String("k", "", "key")
	// RateLimit defines maximum number of concurrent requests to the server.
	// Can be set via flag "-l" or env var "RATE_LIMIT".
	RateLimit = flag.Int("l", 1, "rate limit")

	Loaded = false
)

// ConfigAgent parses command-line flags and environment variables
// to configure agent behavior at runtime.
//
// It respects the following precedence:
// 1. Command-line flags override environment variables.
// 2. Environment variables override defaults.
func ConfigAgent() error {
	var CryptoKeyVar string
	var ConfigPathVar string
	flag.StringVar(&CryptoKeyVar, "crypto-key", "cert.pem", "crypto key")
	flag.StringVar(&ConfigPathVar, "config", "", "config")
	flag.Parse()

	if CryptoKeyVar != "" {
		*CryptoKey = CryptoKeyVar
	}
	if ConfigPathVar != "" {
		*ConfigPath = ConfigPathVar
	}

	// ConfigPath := utils.GetFlagValue("config")
	// if ConfigPath == "" {
	// 	ConfigPath = utils.GetFlagValue("c")
	// }

	if *ConfigPath == "" {
		*ConfigPath = os.Getenv("CONFIG")
	}

	err := LoadConfigFile()
	if err != nil {
		return err
	}
	loadFromEnv()

	// Print current config values
	println("EndpointAgent=", *EndpointAgent)
	println("ReportInterval=", *ReportInterval)
	println("PollInterval=", *PollInterval)
	println("Key=", *Key)
	println("RateLimit=", *RateLimit)
	return nil
}

type Config struct {
	Address        string `json:"address,omitempty"`
	CryptoKey      string `json:"crypto_key,omitempty"`
	ReportInterval string `json:"report_interval,omitempty"`
	PollInterval   string `json:"poll_interval,omitempty"`
}

func loadFromEnv() {
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
	cr, found := os.LookupEnv("CRYPTO_KEY")
	if found {
		CryptoKey = &cr
	}
}

func LoadConfigFile() error {
	cfg := &Config{}
	var a string
	var b int
	var c int
	var d string
	ab, bb, cb, db := false, false, false, false
	if *EndpointAgent != "localhost:8080" {
		a = *EndpointAgent
		ab = true
	}
	if *ReportInterval != 10 {
		b = *ReportInterval
		bb = true
	}
	if *PollInterval != 2 {
		c = *PollInterval
		cb = true
	}
	if *CryptoKey != "cert.pem" {
		d = *CryptoKey
		db = true
	}

	if *ConfigPath != "" {
		file, err := os.Open(*ConfigPath)
		if err != nil {
			return err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err = decoder.Decode(cfg); err != nil {
			return err
		}
		*EndpointAgent = cfg.Address
		dur, err := time.ParseDuration(cfg.ReportInterval)
		if err != nil {
			return err
		}
		*ReportInterval = int(dur.Seconds())
		dur, err = time.ParseDuration(cfg.PollInterval)
		if err != nil {
			return err
		}
		*PollInterval = int(dur.Seconds())
		*CryptoKey = cfg.CryptoKey
		Loaded = true
	}
	checkLoaded(ab, bb, cb, db, a, b, c, d)

	return nil
}

func checkLoaded(ab, bb, cb, db bool, a string, b int, c int, d string) {
	if Loaded {
		if ab {
			*EndpointAgent = a
		}
		if bb {
			*ReportInterval = b
		}
		if cb {
			*PollInterval = c
		}
		if db {
			*CryptoKey = d
		}
		// *EndpointAgent, *ReportInterval , *PollInterval, *CryptoKey = a, b, c, d
	}
}
