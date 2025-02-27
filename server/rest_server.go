package server

import (
	"fmt"
	"log"

	"github.com/Lachstec/digsinet-ng/config"
	"github.com/rs/zerolog/log"
)

func InitRESTServer(cfg config.Configuration) {
	r := NewRESTRouter(cfg)
	err := r.Run(fmt.Sprint(cfg.Port))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to start REST server")
		return
	}
}
