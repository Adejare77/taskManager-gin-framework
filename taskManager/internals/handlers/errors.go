package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "Bad Request",
	})
}

func BadRequestWithMsg(ctx *gin.Context, msg interface{}) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "Bad Request",
		"msg":   msg,
	})
}

func InternalServerError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal Server Error",
	})
}

func InternalServerErrorWithMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal Server Error",
		"msg":   msg,
	})
}

func Unauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
}

func UnauthorizedWithMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
		"msg":   msg,
	})
}

func PageNotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Page Not Found",
		"Available routes": []string{
			"/register",
			"/loging",
			"/tasks",
			"/task/:id",
		},
	})
}
