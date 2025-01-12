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

// func GetTasksByUserID(UserID uint, filter schemas.Task) ([]schemas.Task, error) {
// 	var tasks []schemas.Task

// 	if err := db.Debug().Where("\"userID\" = ? AND title ILIKE '%groceri%' AND status ILIKE '%%' AND \"taskID\" ILIKE '%%'",
// 		UserID).Find(&tasks).Error; err != nil {
// 		return nil, err
// 	}
// 	return tasks, nil
// }

// func GetTasksByTaskID(taskID string) ([]schemas.Task, error) {
// 	var tasks []schemas.Task

// 	db.Scopes(UserObject())

// 	if err := db.First(&tasks, "\"taskID\" = ?", taskID).Error; err != nil {
// 		return nil, err
// 	}
// 	return tasks, nil
// }

// func GetTasksByStatus(status int) ([]schemas.Task, error) {
// 	var tasks []schemas.Task

// 	if err := db.First(&tasks, "status = ?", status).Error; err != nil {
// 		return nil, err
// 	}
// 	return tasks, nil
// }

func GetTaskByTaskID(userID uint, taskID string) schemas.Task {
	var task schemas.Task
	db.First(&task, "\"taskID\" = ? AND \"userID\" = ?", taskID, userID)
	return task
}

func CreateTask(task Task) error {
	if err := db.Create(&task).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTaskByTaskID(userID uint, taskID string) error {
	if err := db.Delete(&Task{}, "\"taskID\" = ? AND \"userID\" = ?", taskID, userID).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTaskByTaskID(userID uint, values schemas.Task) error {
	// if err := db.Scopes(UserObject(userID)).Updates(values).Error; err
	if err := db.Model(&Task{}).Where("taskID = ? AND \"userID\" = ?", values.TaskID, userID).Updates(values).Error; err != nil {
		return err
	}
	return nil
}
