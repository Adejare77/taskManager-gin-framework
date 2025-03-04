package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type APIError struct {
	StatusCode int         `json:"-"`
	ErrorCode  string      `json:"error_code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

func HandleError(ctx *gin.Context, statusCode int, errorCode string, message string, reason string, details interface{}) {
	// For logging
	logrus.WithFields(logrus.Fields{
		"user_id":     ctx.MustGet("userID"),
		"task_id":     ctx.Param("taskID"),
		"status_code": statusCode,
		"error_code":  errorCode,
		"details":     details,
	}).Error(reason)

	// For response to user
	ctx.JSON(statusCode, APIError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
		// Details:    details,
	})
}

func BadRequest(ctx *gin.Context, reason string, details interface{}) {
	HandleError(ctx, http.StatusBadRequest, "BAD_REQUEST", "Invalid request", reason, details)
}

func Unauthorized(ctx *gin.Context, reason string, details interface{}) {
	HandleError(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", reason, details)
}

func InternalServerError(ctx *gin.Context, reason string, details interface{}) {
	HandleError(ctx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An internal server error occured", reason, nil)
}

func NotFound(ctx *gin.Context, reason string, details interface{}) {
	HandleError(ctx, http.StatusOK, "NOT_FOUND", "Task not found", reason, details)
}

func PageNotFound(ctx *gin.Context) {
	HandleError(ctx, http.StatusNotFound, "PAGE_NOT_FOUND", "Page not found", "Invalid Route", gin.H{
		"available_routes": []string{
			"/register",
			"/login",
			"/tasks",
			"/task/:id",
		},
	})
}
