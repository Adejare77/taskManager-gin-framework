package models

import (
	"errors"
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

func GetTaskByTaskID(userID uint, taskID string) (schemas.TaskOutput, error) {
	var task schemas.Task
	cursor := db.Debug().First(&task, "\"taskID\" = ? AND \"userID\" = ?", taskID, userID)

	if cursor.RowsAffected == 0 {
		return schemas.TaskOutput{}, errors.New("task not found")
	}
	return schemas.ToTaskOutput(task), nil
}

func CreateTask(task schemas.Task) error {
	if err := db.Debug().Create(&task).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTaskByTaskID(userID uint, taskID string) error {
	if err := db.Debug().Delete(&schemas.Task{}, "\"taskID\" = ? AND \"userID\" = ?", taskID, userID).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTaskByTaskID(userID uint, taskID string, data map[string]interface{}) error {
	fmt.Println("TASK UPDATING CALLED")
	if err := db.Model(&schemas.Task{}).
		Where("\"taskID\" = ? AND \"userID\" = ?", taskID, userID).Updates(data).Error; err != nil {
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

	if err := db.Model(&schemas.Task{}).Debug().
		Where("\"dueDate\" <= ? AND status != ? AND status != ?", time.Now(), "overdue", "completed").
		Update("status", "overdue").Error; err != nil {
		fmt.Println("Error: ", err)
	}

	if err := db.Model(&schemas.Task{}).Debug().
		Where("\"startDate\" >= ? AND status != ? AND status != ?", time.Now(), "pending", "completed").
		Update("status", "pending").Error; err != nil {
		fmt.Println("Error: ", err)
	}

}
