package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		// Retrieve current user session
		user := session.Get("user")

		if user == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			// Abort Current Execution
			ctx.Abort()
			return
		}

		var userID uint

		// Unmarshal the stored session
		json.Unmarshal(user.([]byte), &userID)

		// store userID in the context for next handler
		ctx.Set("userID", userID)

		// Reset the TTL of key as long as the user is using it
		session.Options(sessions.Options{
			MaxAge: 600,
		})

		session.Save()

		ctx.Next()
	}
}
