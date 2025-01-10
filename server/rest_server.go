package server

// adapted from https://github.com/vsouza/go-gin-boilerplate
import (
	"github.com/Lachstec/digsinet-ng/config"
	"log"
)

func InitRESTServer() {
	conf := config.GetConfig()
	r := NewRESTRouter()
	err := r.Run(conf.GetString("server.port"))
	if err != nil {
		log.Fatal("Unable to start REST server.")
		return
	}
}
