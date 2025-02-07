package configServer

import (
	"flag"
	"os"
)

var (
	Endpoint = flag.String("a", "localhost:8080", "endpoint")
)

func ConfigServer() {
	flag.Parse()
	address, found := os.LookupEnv("ADDRESS")
	if found {
		Endpoint = &address
	}
}
