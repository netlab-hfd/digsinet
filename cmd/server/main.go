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

	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	server.InitRESTServer()
}
