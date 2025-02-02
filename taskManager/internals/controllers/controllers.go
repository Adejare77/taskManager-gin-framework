package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Adejare77/go/taskManager/config"
	"github.com/Adejare77/go/taskManager/internals/handlers"
	"github.com/Adejare77/go/taskManager/internals/models"
	"github.com/Adejare77/go/taskManager/internals/schemas"
	"github.com/Adejare77/go/taskManager/internals/utilities"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func Register(ctx *gin.Context) {
	var user schemas.User
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

	// Create User
	if err := models.Create(user); err != nil {
		fmt.Println("Error:", err.Error())
		if strings.Contains(err.Error(), "duplicate") {
			ctx.JSON(http.StatusConflict, "Email already in use")
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
		fmt.Println("Error: ", err)
		handlers.UnauthorizedWithMsg(ctx, "Invalid email or password")
		return
	}

	if err := utilities.ComparePaswword(user.Password, passwd); err != nil {
		fmt.Println("Error: ", err)
		handlers.UnauthorizedWithMsg(ctx, "Invalid email or password")
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
	title := "%" + ctx.Query("title") + "%"
	status := "%" + ctx.Query("status") + "%"

	var filter schemas.Task
	filter.Title = title
	filter.Status = status

	tasks, _ := models.GetTasksByUserID(userID, filter)

	if len(tasks) == 0 {
		ctx.JSON(http.StatusOK, "Empty")
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

func GetTasksByID(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	taskID := ctx.Param("taskID")

	result, err := models.GetTaskByTaskID(userID, taskID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func DeleteTask(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	taskID := ctx.Param("taskID")

	if err := models.DeleteTaskByTaskID(userID, taskID); err != nil {
		fmt.Println("Error:", err)
		handlers.InternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, "Task successfully deleted")
	ctx.Redirect(http.StatusSeeOther, "/task")
}

func PostTask(ctx *gin.Context) {
	var task schemas.PostTask
	userID := ctx.MustGet("userID").(uint)

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

	body, err := schemas.ToTask(task)
	if err != nil {
		fmt.Println("Error: ", err)
		handlers.BadRequest(ctx)
		return
	}

	body.UserID = userID
	body.TaskID = uuid.New().String()

	if err := models.CreateTask(body); err != nil {
		fmt.Println("Error:", err)
		handlers.BadRequestWithMsg(ctx, err.Error())
		return
	}

	result, err := models.GetTaskByTaskID(userID, body.TaskID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func UpdateTask(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	taskID := ctx.Param("taskID")

	var dataValues map[string]interface{}
	var dataTime schemas.DateTimeUpdate

	if err := ctx.ShouldBindBodyWith(&dataValues, binding.JSON); err != nil {
		fmt.Println("Error: ", err)
		handlers.BadRequest(ctx)
		return
	}

	if err := ctx.ShouldBindBodyWith(&dataTime, binding.JSON); err != nil {
		if fieldErros, ok := err.(validator.ValidationErrors); ok {
			msg := utilities.ValidationError(fieldErros)
			handlers.BadRequestWithMsg(ctx, msg)
			fmt.Println("Error: ", err)
			return
		}
	}

	if dataTime.StartDate != "" || dataTime.DueDate != "" {
		existingDoc, err := models.GetTaskByTaskID(userID, taskID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		if dataTime.StartDate == "" {
			dataTime.StartDate = existingDoc.StartDate
		}
		if dataTime.DueDate == "" {
			dataTime.DueDate = existingDoc.DueDate
		}

		var startDate time.Time
		if startDate, err = utilities.StartTimeManipulator(dataTime.StartDate); err == nil {
			dataValues["startDate"] = startDate
		} else {
			fmt.Println("Error: ", err)
			handlers.BadRequestWithMsg(ctx, err.Error())
			return
		}

		if dueDate, err := utilities.DueTimeManipulator(dataTime.DueDate, startDate); err == nil {
			dataValues["dueDate"] = dueDate
		} else {
			fmt.Println("Error: ", err)
			handlers.BadRequestWithMsg(ctx, err.Error())
			return
		}
	}

	if err := models.UpdateTaskByTaskID(userID, taskID, dataValues); err != nil {
		fmt.Println("Error: ", err)
		handlers.BadRequestWithMsg(ctx, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, fmt.Sprint(taskID))
}

func DeleteUser(ctx *gin.Context) {
	UserID := ctx.MustGet("userID").(string)

	if err := models.DeleteUser(UserID); err != nil {
		fmt.Println(err)
		handlers.InternalServerErrorWithMsg(ctx, "Unable to delete User")
		return
	}

	config.DeleteSession(ctx)

	ctx.JSON(http.StatusOK, "User Deleted Successfully")
}
