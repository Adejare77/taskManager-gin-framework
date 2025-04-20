package models

import (
	"errors"
	"time"

	"github.com/adejare77/taskmanager-gin-framework/config"
	"github.com/adejare77/taskmanager-gin-framework/internals/schemas"
	"gorm.io/gorm"
)

func CreateTask(task schemas.Task) error {
	return config.DB.Create(&task).Error
}

func FindTasksByUserID(userID string, filter schemas.TaskQueryParams) ([]map[string]any, error) {
	var tasks []map[string]any
	offset := (filter.Page - 1) * filter.Limit

	query := config.DB.Model(&schemas.Task{}).
		Where("title ILIKE ? AND status ILIKE ? AND user_id = ?",
			filter.Title, filter.Status, userID).
		Select(`
		task_id,
		"desc" AS description,
		title,
		start_date,
		due_date,
		status
		`).
		Order("created_at DESC").
		Offset(offset).
		Limit(filter.Limit)

	if query.Error != nil {
		return nil, query.Error
	}

	query.Scan(&tasks)

	return tasks, nil
}

func FindTaskByTaskID(userID string, taskID string) (map[string]any, error) {
	query := config.DB.Model(&schemas.Task{}).
		Where("task_id = ? AND user_id = ?", taskID, userID).
		Select(`
		task_id,
		"desc" AS description,
		title,
		start_date,
		due_date,
		status
		`)

	if query.Error != nil {
		return nil, query.Error
	}

	var task map[string]any
	query.Scan(&task)

	return task, nil
}

func UpdateTask(userID string, taskID string, data map[string]any) error {
	cursor := config.DB.Model(&schemas.Task{}).
		Where("user_id = ? AND task_id = ?", userID, taskID).
		Updates(data)

	if cursor.Error != nil {
		return cursor.Error
	}

	if cursor.RowsAffected == 0 {
		return errors.New("not found or unauthorized")
	}
	return nil
}

func DeleteTask(userID string, taskID string) error {
	cursor := config.DB.Delete(&schemas.Task{}, "task_id = ? AND user_id = ?",
		taskID, userID)

	if cursor.Error != nil {
		return cursor.Error
	}

	if cursor.RowsAffected == 0 {
		return errors.New("not found or unauthorized")
	}
	return nil
}

func StatusUpdater() {
	config.DB.Model(&schemas.Task{}).
		Where("status IN ?", []string{"pending", "in-progress"}).
		Update(
			"status",
			gorm.Expr(`
			CASE
			WHEN status = 'in-progress' AND due_date < ? THEN 'overdue'
			WHEN status = 'pending' AND start_date < ? THEN 'in-progress'
			ELSE status
			END
			`, time.Now(), time.Now()),
		)
}
