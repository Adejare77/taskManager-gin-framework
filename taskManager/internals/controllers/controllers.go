package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Adejare77/go/taskManager/config"
	"github.com/Adejare77/go/taskManager/internals/handlers"
	"github.com/Adejare77/go/taskManager/internals/models"
	"github.com/Adejare77/go/taskManager/internals/schemas"
	"github.com/Adejare77/go/taskManager/internals/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Println("Error:", err)
		if fieldError, ok := err.(validator.ValidationErrors); ok {
			msg := utilities.ValidationError(fieldError)
			handlers.BadRequestWithMsg(ctx, msg)
			return
		}
		handlers.BadRequest(ctx)
		return
	}

	// hash password
	hashedPwd, err := utilities.HashPassword(user.Password)
	if err != nil {
		fmt.Println("Error:", err)
		handlers.InternalServerError(ctx)
		return
	}

	// Set hashedPwd value
	user.Password = string(hashedPwd)

	// Create User
	if err := user.Create(); err != nil {
		fmt.Println("Error:", err.Error())
		if strings.Contains(err.Error(), "duplicate") {
			handlers.BadRequestWithMsg(ctx, "Email Already Exists")
		} else {
			handlers.InternalServerErrorWithMsg(ctx, "Could Not Register User")
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"fullName": user.FullName,
		"email":    user.Email,
		"status":   "Successfully Registered",
	})
}

func Login(ctx *gin.Context) {
	var user schemas.Login
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Println("Error:", err)
		if fieldError, ok := err.(validator.ValidationErrors); ok {
			msg := utilities.ValidationError(fieldError)
			handlers.BadRequestWithMsg(ctx, msg)
			return
		}
		handlers.BadRequest(ctx)
		return
	}

	userID, passwd, err := models.GetInfo(user.Email)
	if err != nil {
		handlers.UnauthorizedWithMsg(ctx, "Invalid Email or Password")
		return
	}

	if err := utilities.ComparePaswword(user.Password, passwd); err != nil {
		handlers.UnauthorizedWithMsg(ctx, "Invalid Email or Password")
		return
	}

	if err := config.CreateSession(ctx, userID); err != nil {
		fmt.Println("Error:", err)
		handlers.InternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, "login Successfully")
}

func GetTasks(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	title := ctx.Query("title")
	status := ctx.Query("status")
	taskID := ctx.Query("taskID")

	if title == "" {
		title = "%"
	}
	if status == "" {
		status = "%"
	}
	if taskID == "" {
		taskID = "%"
	}

	var filter schemas.Task
	filter.Title = title
	filter.Status = status
	filter.TaskID = taskID

	// fmt.Println("-------------------------------")
	// fmt.Println(status)
	// fmt.Println(title)
	// fmt.Println(taskID)
	// fmt.Println("-------------------------------")

	tasks, err := models.GetTasksByUserID(userID, filter)
	if err != nil {
		fmt.Println(tasks)
		fmt.Println(err)
		if strings.Contains(err.Error(), "record not found") {
			ctx.JSON(http.StatusOK, "Empty Task. Add New Task")
			return
		}
		handlers.InternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func GetTasksByID(ctx *gin.Context) {
	ctx.MustGet("userID")
	taskID := ctx.Param("taskID")

	result, err := models.GetTasksByTaskID(taskID)
	if err != nil {
		fmt.Println("Error: ", err)
		handlers.InternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func DeleteTask(ctx *gin.Context) {
	ctx.MustGet("userID")

	type TaskID struct {
		TaskID string `json:"taskID" binding:"required"`
	}
	var taskID TaskID

	if err := ctx.ShouldBindJSON(&taskID); err != nil {
		fmt.Println("Error:", err)
		handlers.BadRequestWithMsg(ctx, "`id` field required")
		return
	}

	if err := models.DeleteTaskByTaskID(taskID.TaskID); err != nil {
		fmt.Println("Error:", err)
		handlers.InternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, "Task Deleted Successfully")
	ctx.Redirect(http.StatusSeeOther, "/task")
}

func PostTask(ctx *gin.Context) {
	var task models.Task
	task.UserID = ctx.MustGet("userID").(uint)

	if err := ctx.ShouldBindBodyWithJSON(&task); err != nil {
		fmt.Println("Error:", err)
		if fieldErrors, ok := err.(validator.ValidationErrors); ok {
			msg := utilities.ValidationError(fieldErrors)
			handlers.BadRequestWithMsg(ctx, msg)
			return
		}
		handlers.BadRequest(ctx)
		return
	}

	// inputDate, _ := task.DueDate.(string)
	task.DueDate = utilities.TimeManipulator(task.DueDate)

	// Assign a Unique Task ID
	task.TaskID = uuid.New().String()

	if err := models.CreateTask(task); err != nil {
		fmt.Println("Error:", err)
		handlers.InternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"taskID":      task.TaskID,
		"title":       task.Title,
		"description": task.Desc,
		"dueDate":     task.DueDate,
		"status":      task.Status,
	})
}

func UpdateTask(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	var task schemas.Task

	if err := ctx.ShouldBindBodyWithJSON(&task); err != nil {
		fmt.Println("Error:", err)
		handlers.BadRequest(ctx)
		return
	}

	if err := models.UpdateTaskByUserID(userID, task); err != nil {
		fmt.Println("Error:", err)
		handlers.InternalServerErrorWithMsg(ctx, "Error Updating Task")
		return
	}

	ctx.JSON(http.StatusOK, task)
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/task/:%v", userID))
}

func DeleteUser(ctx *gin.Context) {
	UserID := ctx.MustGet("userID").(string)

	if err := models.DeleteUser(UserID); err != nil {
		fmt.Println(err)
		handlers.InternalServerErrorWithMsg(ctx, "Unable to delete User")
		return
	}

	// Delete the User's Session
	config.DeleteSession(ctx)

	ctx.JSON(http.StatusOK, "User Deleted Successfully")
}
