package middlewares

import (
	"github.com/Adejare77/go/taskManager/internals/handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		// Retrieve current user session
		user := session.Get("currentUser")

		if user == nil {
			handlers.Unauthorized(ctx, "login required", "unauthorized access")
			ctx.Abort()
			return
		}

		// store userID in the context for next handler
		ctx.Set("currentUser", user)     // set key for next function
		session.Set("currentUser", user) // roll-over key
		session.Save()

		ctx.Next()
	}
}
