package handlers

import (
	"net/http"

	"github.com/Adejare77/go/taskManager/internals/utilities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type APIError struct {
	StatusCode int    `json:"-"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
}

func Info(details map[string]any, message string) {
	logrus.WithFields(logrus.Fields(details)).Info(message)

}

func Warning(message string) {
	logrus.Warning(message)

}

func HandleError(ctx *gin.Context, statusCode int, errorCode string, message any, details any) {
	// For logging
	logrus.WithFields(logrus.Fields{
		"status_code": statusCode,
		"error_code":  errorCode,
		"details":     details,
	}).Error(details)

	// For response to user
	ctx.JSON(statusCode, gin.H{
		"status": statusCode,
		"error":  message,
	})
}

func BadRequest(ctx *gin.Context, message string, details any) {
	HandleError(ctx, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

func Unauthorized(ctx *gin.Context, message string, details any) {
	HandleError(ctx, http.StatusUnauthorized, "UNAUTHORIZED", message, details)
}

func InternalServerError(ctx *gin.Context, message string, details any) {
	HandleError(ctx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, details)
}

func NotFound(ctx *gin.Context, message string, details any) {
	HandleError(ctx, http.StatusOK, "NOT_FOUND", message, details)
}

func Validation(ctx *gin.Context, err error) {
	validationDetails := utilities.ValidationError(err)
	HandleError(ctx, http.StatusBadRequest, "VALIDATION_ERROR", validationDetails, err)
}

// func PageNotFound(ctx *gin.Context) {
// 	HandleError(ctx, http.StatusNotFound, "PAGE_NOT_FOUND", "Page not found", "Invalid Route", gin.H{
// 		"available_routes": []string{
// 			"/register",
// 			"/login",
// 			"/tasks",
// 			"/task/:id",
// 		},
// 	})
// }
