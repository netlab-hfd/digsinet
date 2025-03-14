package rest_middlewares

// adapted from https://github.com/vsouza/go-gin-boilerplate

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"

	"crypto/sha256"
	"crypto/subtle"

	"github.com/Lachstec/digsinet-ng/config"
	"github.com/gin-gonic/gin"
)

func sha256Sum(s string) []byte {
	sum := sha256.Sum256([]byte(s))
	arr := make([]byte, len(sum))
	copy(arr, sum[:])

	return arr
}

// secureCompare calculates sha256 hash of parameters a and b and does constant time comparison
// to avoid time based attacks.
func secureCompare(a, b string) int {
	aSum := sha256Sum(a)
	bSum := sha256Sum(b)

	return subtle.ConstantTimeCompare(aSum, bSum)
}

func AuthMiddleware(cfg config.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqKey := c.Request.Header.Get("X-Auth-Key")
		reqSecret := c.Request.Header.Get("X-Auth-Secret")

		var key string
		var secret string
		if key = cfg.Http.AuthKey; len(strings.TrimSpace(key)) == 0 {
			err := c.AbortWithError(500, errors.New("no authentication key provided"))
			if err != nil {
				log.Warn().
					Err(err).
					Msg("Missing Authentication Key")
			}
			return
		}
		if secret = cfg.Http.AuthSecret; len(strings.TrimSpace(secret)) == 0 {
			err := c.AbortWithError(401, errors.New("no authentication secret provided"))
			if err != nil {
				log.Warn().
					Err(err).
					Msg("Missing Authentication Secret")
			}
			return
		}

		isKeysEqual := secureCompare(key, reqKey) == 1
		isSecretsEqual := secureCompare(secret, reqSecret) == 1
		if !isKeysEqual || !isSecretsEqual {
			err := c.AbortWithError(401, errors.New("authentication failed"))
			if err != nil {
				log.Warn().
					Err(err).
					Msg("Authentication Failed")
			}
			return
		}
		c.Next()
	}
}
