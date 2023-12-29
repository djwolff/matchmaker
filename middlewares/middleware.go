package middlewares

import (
	"net/http"

	"github.com/djwolff/matchmaker/utils/token"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := token.ExtractUserFromToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Attach the user to the context for further processing
		c.Set("user", user)

		c.Next()
	}
}
