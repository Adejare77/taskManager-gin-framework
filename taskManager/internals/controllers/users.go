package controllers

import (
	"net/http"
	"strings"

	"github.com/Adejare77/go/taskManager/config"
	"github.com/Adejare77/go/taskManager/internals/handlers"
	"github.com/Adejare77/go/taskManager/internals/models"
	"github.com/Adejare77/go/taskManager/internals/schemas"
	"github.com/Adejare77/go/taskManager/internals/utilities"
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user schemas.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		handlers.Validation(ctx, err)
		return
	}

	if err := models.Create(user); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			handlers.BadRequest(ctx, "email already exists", err)
			return
		}
		handlers.InternalServerError(ctx, "internal error", err)
		return
	}

	handlers.Info(gin.H{"email": user.Email}, "New User Registered")

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "registered successfully",
	})
}

func Login(ctx *gin.Context) {
	type Login struct {
		Email    string `binding:"required"`
		Password string `binding:"required"`
	}

	var user Login
	if err := ctx.ShouldBindJSON(&user); err != nil {
		handlers.Validation(ctx, err)
		return
	}

	currentUser, err := models.FindUserInfo(user.Email)
	if err != nil {
		handlers.Unauthorized(ctx, "invalid email or password", err)
		return
	}

	if err := utilities.ComparePaswword(user.Password, currentUser.Password); err != nil {
		handlers.Unauthorized(ctx, "invalid email or password", err)
		return
	}

	if err := config.CreateSession(ctx, currentUser.ID); err != nil {
		handlers.InternalServerError(ctx, "failed to create session", err)
		return
	}

	handlers.Info(gin.H{"user_id": currentUser.Email}, "user logged in")

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "login successful",
	})

}

func Logout(ctx *gin.Context) {
	_, exists := ctx.Get("currentUser")
	if !exists {
		handlers.Unauthorized(ctx, "user needs to log in", "Logout Error")
		return
	}

	config.DeleteSession(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "logout successful",
	})
}

func DeleteUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(string)

	if err := models.DeleteUser(currentUser); err != nil {
		handlers.InternalServerError(ctx, "could not delete user", err)
		return
	}

	config.DeleteSession(ctx)

	ctx.JSON(http.StatusOK, "User deleted successfully")
}
