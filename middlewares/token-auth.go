package middlewares

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/majori/bolley/controllers"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			controllers.RespondWithError(401, "API token required", c)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") ||
			(parts[1] != os.Getenv("API_TOKEN")) {
			controllers.RespondWithError(401, "Invalid API token", c)
			return
		}

		c.Next()
	}
}
