package models

import (
	"fmt"
	"time"

	"github.com/Adejare77/go/taskManager/internals/schemas"
)

func GetTasksByUserID(userID uint, filter schemas.Task) ([]schemas.TaskOutput, error) {
	var tasks []schemas.Task

	if err := db.Debug().Where("title ILIKE ? AND status ILIKE ? AND \"userID\" = ?",
		filter.Title, filter.Status, userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	var taskOutput []schemas.TaskOutput
	for _, task := range tasks {
		taskOutput = append(taskOutput, schemas.ToTaskOutput(task))
	}
	return taskOutput, nil
}

func GetTaskByTaskID(userID uint, taskID string) schemas.TaskOutput {
	var task schemas.Task
	db.Debug().First(&task, "\"taskID\" = ? AND \"userID\" = ?", taskID, userID)
	return schemas.ToTaskOutput(task)
}

func CreateTask(task schemas.Task) error {
	if err := db.Debug().Create(&task).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTaskByTaskID(userID uint, taskID string) error {
	if err := db.Debug().Delete(&schemas.TaskOutput{}, "\"taskID\" = ? AND \"userID\" = ?", taskID, userID).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTaskByTaskID(userID uint, values schemas.Task) error {
	if err := db.Model(&schemas.TaskOutput{}).Where("taskID = ? AND \"userID\" = ?", values.TaskID, userID).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

func CheckStatus() {
	if err := db.Model(&schemas.Task{}).Debug().
		Where("\"startDate\" <= ? AND \"dueDate\" > ? AND status != ?", time.Now(), time.Now(), "in-progress").
		Update("status", "in-progress").Error; err != nil {
		fmt.Println("Error: ", err)
	}

	db.Model(&schemas.Task{}).Debug().
		Where("\"dueDate\" <= ? AND status != ? AND status != ?", time.Now(), "overdue", "completed").
		Update("status", "overdue")
}
