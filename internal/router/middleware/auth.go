package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/gin-gonic/gin"
	"net/http"
	"start/internal/config"
)

func BasicAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(cfg.Username))
			expectedPasswordHash := sha256.Sum256([]byte(cfg.Password))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				c.Next()
				return
			}
		}

		c.Writer.Header().Set(
			"WWW-Authenticate",
			`Basic realm="restricted", charset="UTF-8"`,
		)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
	}
}
