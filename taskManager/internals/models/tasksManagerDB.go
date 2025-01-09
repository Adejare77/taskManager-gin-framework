package models

import (
	"time"

	"github.com/Adejare77/go/taskManager/internals/schemas"
)

type Task struct {
	UserID    uint   `gorm:"column:userID;index;not null"`
	TaskID    string `gorm:"column:taskID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Desc      string `json:"description" gorm:"column:description;not null" binding:"required"`
	Title     string `json:"title" gorm:"column:title;not null" binding:"required"`
	DueDate   string `json:"dueDate" gorm:"column:dueDate;not null" binding:"required,dueDate"`
	Status    string `json:"status" gorm:"column:status;not null" binding:"required,status"`
}

func GetTasksByUserID(UserID uint, filter schemas.Task) ([]schemas.Task, error) {
	var tasks []schemas.Task

	if err := db.Where("\"userID\" = ? AND title ILIKE ? AND status ILIKE ? AND \"taskID\" ILIKE ?",
		UserID, filter.Title, filter.Status, filter.TaskID).Find(&Task{}).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByTaskID(taskID string) ([]schemas.Task, error) {
	var tasks []schemas.Task

	if err := db.First(&tasks, "\"taskID\" = ?", taskID).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByStatus(status int) ([]schemas.Task, error) {
	var tasks []schemas.Task

	if err := db.First(&tasks, "status = ?", status).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func CreateTask(task Task) error {
	if err := db.Create(&task).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTaskByTaskID(taskID string) error {
	if err := db.Delete(&Task{}, "taskID = ?", taskID).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTaskByUserID(userID uint, values schemas.Task) error {
	if err := db.Model(&Task{}).Updates(values).Error; err != nil {
		return err
	}
	return nil
}
