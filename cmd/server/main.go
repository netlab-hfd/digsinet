package main

// adapted from https://github.com/vsouza/go-gin-boilerplate

import (
	"flag"
	"fmt"
	"os"

	"github.com/Lachstec/digsinet-ng/log"

	"github.com/Lachstec/digsinet-ng/config"
	"github.com/Lachstec/digsinet-ng/server"
)

//nolint:all
func main() {
	log.InitLogging()

	var eFlag = flag.String("e", "", "Environment mode (determining the config file to look for)")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	config.Init(*eFlag)
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("failed to read in configuration file")
	}

	server.InitRESTServer(*cfg)
}
