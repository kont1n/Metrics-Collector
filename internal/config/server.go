package config

import (
	"flag"
	"os"
)

func parseServerFlags() {
	flag.StringVar(&flagRunAddr, "a", DefaultServer, "address and port to run server")
	flag.Parse()
}

func ParseServerConfig() (host string) {
	envAddress := os.Getenv("ADDRESS")
	if envAddress != "" {
		host = envAddress
	} else {
		parseServerFlags()
		host = flagRunAddr
	}
	return host
}
