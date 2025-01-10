package server

// adapted from https://github.com/vsouza/go-gin-boilerplate
import (
	"github.com/Lachstec/digsinet-ng/config"
	"github.com/rs/zerolog/log"
)

func InitRESTServer() {
	conf := config.GetConfig()
	r := NewRESTRouter()
	err := r.Run(conf.GetString("server.port"))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to start REST server")
		return
	}
}
