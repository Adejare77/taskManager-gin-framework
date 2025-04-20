package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Adejare77/taskmanager/internals/handlers"
	"github.com/Adejare77/taskmanager/internals/models"
	"github.com/Adejare77/taskmanager/internals/schemas"
	"github.com/Adejare77/taskmanager/internals/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func PostTask(ctx *gin.Context) {
	var task schemas.Task

	// default values
	task.UserID = ctx.MustGet("currentUser").(string)

	if err := ctx.ShouldBindJSON(&task); err != nil {
		if _, check := err.(validator.ValidationErrors); check {
			handlers.Validation(ctx, err)
		} else {
			handlers.BadRequest(ctx,
				"validation failed. invalid date format. Use `YYYY-MM-DAY HH:MM` e.g 2025-03-21 00:05",
				"Invalid date field format",
			)
		}
		return
	}

	if startDate, err := utilities.CompareDates(task.StartDate, task.DueDate); err != nil {
		handlers.BadRequest(ctx, err.Error(), err)
		return
	} else {
		if task.StartDate == nil {
			task.StartDate = new(utilities.JSONTime)
		}
		*task.StartDate = startDate
	}

	if err := models.CreateTask(task); err != nil {
		handlers.InternalServerError(ctx, "could not create new task", err)
		return
	}

	handlers.Info(map[string]any{
		"user_id": task.UserID,
		"title":   task.Title,
	}, "created new task")

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "task created successful",
		"message": http.StatusOK,
	})
}

func GetTasks(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(string)

	var filters schemas.TaskQueryParams

	// Set Default Value
	filters.Page = 1
	filters.Limit = 20

	if err := ctx.ShouldBindQuery(&filters); err != nil {
		handlers.Validation(ctx, err)
		return
	}

	filters.Title = "%" + filters.Title + "%"
	filters.Status = "%" + filters.Status + "%"

	tasks, err := models.FindTasksByUserID(userID, filters)
	if err != nil {
		handlers.InternalServerError(ctx, "failed to fetch tasks", err)
		return
	}

	if len(tasks) == 0 {
		ctx.JSON(http.StatusOK, []string{})
		return
	}

	next, prev := generateLink(filters, len(tasks))

	ctx.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"meta": gin.H{
			"page":  filters.Page,
			"limit": filters.Limit,
			"_links": gin.H{
				"prev": prev,
				"next": next,
			},
		},
	})
}

func GetTasksByID(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(string)

	var task schemas.TaskUriParam

	if err := ctx.ShouldBindUri(&task); err != nil {
		handlers.Validation(ctx, err)
		return
	}

	result, err := models.FindTaskByTaskID(userID, task.TaskID)
	if err != nil {
		handlers.NotFound(ctx, "failed to fetch task", err)
		return
	}

	if len(result) == 0 {
		handlers.BadRequest(ctx, "task not found", "taskID not present")
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func UpdateTask(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(string)

	var task schemas.TaskUriParam
	if err := ctx.ShouldBindUri(&task); err != nil {
		handlers.Validation(ctx, err)
		return
	}

	type postUpdateDTO struct {
		Desc      *string             `json:"description" binding:"omitempty"`
		Title     *string             `binding:"omitempty"`
		StartDate *utilities.JSONTime `json:"start_date" binding:"omitempty"`
		DueDate   *utilities.JSONTime `json:"due_date" binding:"omitempty"`
	}

	var dto postUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		handlers.Validation(ctx, err)
		return
	}

	updateData := make(map[string]any)

	if dto.Title != nil {
		updateData["title"] = dto.Title
	}
	if dto.Desc != nil {
		updateData["desc"] = dto.Desc
	}

	if dto.StartDate != nil || dto.DueDate != nil {
		if dto.StartDate == nil || dto.DueDate == nil {
			old_task, err := models.FindTaskByTaskID(userID, task.TaskID)
			if err != nil {
				handlers.InternalServerError(ctx, err.Error(), err)
				return
			}

			if dto.StartDate == nil {
				v := utilities.JSONTime(old_task["start_date"].(time.Time))
				dto.StartDate = &v
			} else {
				v := utilities.JSONTime(old_task["due_date"].(time.Time))
				dto.DueDate = &v
			}
		}

		_, err := utilities.CompareDates(dto.StartDate, *dto.DueDate)
		if err != nil {
			handlers.BadRequest(ctx, err.Error(), err)
			return
		}

		updateData["start_date"] = dto.StartDate
		updateData["due_date"] = dto.DueDate
	}

	if err := models.UpdateTask(userID, task.TaskID, updateData); err != nil {
		if strings.Contains(err.Error(), "not found") {
			handlers.BadRequest(ctx, "not found or unauthorized", err)
			return
		}
		handlers.InternalServerError(ctx, fmt.Sprintf("could not update data: %s", task.TaskID), err)
		return
	}

	handlers.Info(map[string]any{
		"task_id": task.TaskID,
	}, "updated successfully")

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "updated successfully",
	})
}

func DeleteTask(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(string)
	var task schemas.TaskUriParam

	if err := ctx.ShouldBindUri(&task); err != nil {
		handlers.Validation(ctx, err)
		return
	}

	if err := models.DeleteTask(userID, task.TaskID); err != nil {
		handlers.InternalServerError(ctx, fmt.Sprintf("could not delete task: %s", task.TaskID), err)
		return
	}

	handlers.Info(map[string]any{
		"user_id": userID,
		"task_id": task.TaskID,
	}, "task deleted")

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("task %s deleted", task.TaskID),
	})
}

func generateLink(filters schemas.TaskQueryParams, totalTasks int) (string, string) {
	prev := fmt.Sprintf("/tasks?page=%d&limit=%d", filters.Page-1, filters.Limit)
	next := fmt.Sprintf("/tasks?page=%d&limit=%d", filters.Page+1, filters.Limit)

	if filters.Page == 1 {
		prev = "null"
	}

	if filters.Page >= totalTasks+(filters.Limit-1)/filters.Limit {
		next = "null"
	}

	return next, prev
}
