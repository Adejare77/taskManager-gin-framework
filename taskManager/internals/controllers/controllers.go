package controllers

import (
	"log"
	"net/http"

	"github.com/Adejare77/go/taskManager/config"
	"github.com/Adejare77/go/taskManager/internals/handlers"
	"github.com/Adejare77/go/taskManager/internals/models"
	"github.com/Adejare77/go/taskManager/internals/schemas"
	"github.com/Adejare77/go/taskManager/internals/utilities"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Register(ctx *gin.Context) {
	var user schemas.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		handlers.BadRequest(ctx, "Bad Request", utilities.ValidationError(err))
		return
	}

	if err := models.Create(user); err != nil {
		handlers.InternalServerError(ctx, "Failed to create user", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"email": user.Email,
	}).Info("User registered successfully")

	ctx.JSON(http.StatusCreated, gin.H{
		"fullName": user.FullName,
		"email":    user.Email,
		"status":   "Successfully Registered",
	})
}

func Login(ctx *gin.Context) {
	var user schemas.Login
	if err := ctx.ShouldBindJSON(&user); err != nil {
		handlers.BadRequest(ctx, "Invalid login format", utilities.ValidationError(err))
		return
	}

	// log.Println("-----------------------------")
	// log.Println(user.Email)
	// log.Println(user.Password)
	// log.Println("-----------------------------")

	userID, passwd, err := models.GetInfo(user.Email)
	if err != nil {
		handlers.Unauthorized(ctx, "Failed to fetch user info", "Invalid email or password")
		return
	}
	log.Println("-----------------------------")
	log.Println("POINT 2")
	log.Println(user)
	log.Println("-----------------------------")

	if err := utilities.ComparePaswword(user.Password, passwd); err != nil {
		handlers.Unauthorized(ctx, "Invalid password", "Invalid email or password")
		return
	}

	if err := config.CreateSession(ctx, userID); err != nil {
		handlers.InternalServerError(ctx, "Failed to create session", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("User logged in successfully")

	ctx.JSON(http.StatusOK, "Login Successful")
}

func GetTasks(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	title := "%" + ctx.Query("title") + "%"
	status := "%" + ctx.Query("status") + "%"

	var filter schemas.Task
	filter.Title = title
	filter.Status = status

	tasks, err := models.GetTasksByUserID(userID, filter)
	if err != nil {
		handlers.InternalServerError(ctx, "Failed to fetch tasks", err)
		return
	}

	if len(tasks) == 0 {
		ctx.JSON(http.StatusOK, "No tasks found")
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func GetTasksByID(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	taskID := ctx.Param("taskID")

	result, err := models.GetTaskByTaskID(userID, taskID)
	if err != nil {
		handlers.NotFound(ctx, "Failed to fetch task", err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func DeleteTask(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	taskID := ctx.Param("taskID")

	if err := models.DeleteTaskByTaskID(userID, taskID); err != nil {
		handlers.InternalServerError(ctx, "Failed to delete task", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id": userID,
		"task_id": taskID,
	}).Info("Task deleted Successfully")

	ctx.JSON(http.StatusOK, "Task successfully deleted")
}

func PostTask(ctx *gin.Context) {
	var task schemas.PostTask
	userID := ctx.MustGet("userID").(uint)

	if err := ctx.ShouldBindBodyWithJSON(&task); err != nil {
		handlers.BadRequest(ctx, "Invalid Post Request Body", utilities.ValidationError(err))
		return
	}

	body, err := schemas.ToTask(task)
	if err != nil {
		handlers.BadRequest(ctx, "Could not Convert PostTask Input to Task", err.Error())
		return
	}

	body.UserID = userID
	body.TaskID = uuid.New().String()

	if err := models.CreateTask(body); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error()

		handlers.InternalServerError(ctx, "Failed to create task", err)
		return
	}

	result, err := models.GetTaskByTaskID(userID, body.TaskID)
	if err != nil {
		handlers.NotFound(ctx, "Task not found", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"use_id":  userID,
		"task_id": body.TaskID,
	}).Info("Task created successfully")

	ctx.JSON(http.StatusOK, result)
}

func UpdateTask(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	taskID := ctx.Param("taskID")

	var dataValues map[string]interface{}
	var dateData schemas.DateTimeUpdate

	// Bind JSON data to dataValues and dateData
	if err := ctx.ShouldBindBodyWith(&dataValues, binding.JSON); err != nil {
		handlers.BadRequest(ctx, "Invalid Request Body", utilities.ValidationError(err))
		return
	}

	if err := ctx.ShouldBindBodyWith(&dateData, binding.JSON); err != nil {
		handlers.BadRequest(ctx, "Invalid date fields in Request Body", utilities.ValidationError(err))
		return
	}

	if dateData.StartDate != "" || dateData.DueDate != "" {
		// Fetch the existing task to get current dates
		existingTask, err := models.GetTaskByTaskID(userID, taskID)
		if err != nil {
			handlers.NotFound(ctx, "Task not found", err)
		}

		// if startDate is not provided, use the existing startDate
		if dateData.StartDate == "" {
			dateData.StartDate = existingTask.StartDate
		}

		// if dueDate is not provided, use the existing dueDate
		if dateData.DueDate == "" {
			dateData.DueDate = existingTask.DueDate
		}

		// validate and parse startDate
		startDate, err := utilities.StartTimeManipulator(dateData.StartDate)
		if err != nil {
			handlers.BadRequest(ctx, "Invalid startDate", err.Error())
			return
		}

		// validate and parse dueDate
		dueDate, err := utilities.DueTimeManipulator(dateData.DueDate, startDate)
		if err != nil {
			handlers.BadRequest(ctx, "Invalid dueDate", err.Error())
			return
		}

		// Update dataValues with the parsed dates
		dataValues["startDate"] = startDate
		dataValues["dueDate"] = dueDate

		// Update the task with the new data
		if err := models.UpdateTaskByTaskID(userID, taskID, dataValues); err != nil {
			handlers.InternalServerError(ctx, "Unable to update data", err.Error())
			return
		}

		ctx.Redirect(http.StatusSeeOther, taskID)
	}

	if err := models.UpdateTaskByTaskID(userID, taskID, dataValues); err != nil {
		handlers.InternalServerError(ctx, "Failed to update task", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id": userID,
		"task_id": taskID,
	}).Info("Task updated successfully")

	ctx.JSON(http.StatusOK, "Task updated successfully")
}

func DeleteUser(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(int)

	if err := models.DeleteUser(userID); err != nil {
		handlers.InternalServerError(ctx, "Could not delete user", err)
		return
	}

	config.DeleteSession(ctx)

	ctx.JSON(http.StatusOK, "User deleted successfully")
}
