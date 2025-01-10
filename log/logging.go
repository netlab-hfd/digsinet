package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// InitLogging initializes the default logger to use a zerolog
// instance that logs to stdout by default.
func InitLogging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}
