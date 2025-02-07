package config

import (
	"flag"
	"os"
)

var (
	EndpointS = flag.String("a", "localhost:8080", "endpoint")
)

func ConfigServer() {
	flag.Parse()
	address, found := os.LookupEnv("ADDRESS")
	if found {
		EndpointS = &address
	}
}
