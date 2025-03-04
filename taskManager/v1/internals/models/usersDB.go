package models

import (
	"errors"

	"github.com/Adejare77/go/taskManager/config"
	"github.com/Adejare77/go/taskManager/internals/schemas"
	"gorm.io/gorm"
)

var db = config.DB

func GetTasksByStatus(userID uint, status string) ([]schemas.TaskOutput, error) {
	var user schemas.User
	if err := db.Preload("Tasks", "status = ?", status).First(&user, userID).Error; err != nil {
		return nil, err
	}

	var taskOutput []schemas.TaskOutput
	for _, task := range user.Tasks {
		taskOutput = append(taskOutput, schemas.ToTaskOutput(task))
	}
	return taskOutput, nil
}

func GetInfo(email string) (uint, string, error) {
	var user schemas.User
	err := db.Debug().Where("email=?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, "", err
	}
	return user.ID, user.Password, nil
}

func Create(user schemas.User) error {
	err := db.Create(&user).Error

	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(userID string) error {
	if err := db.Unscoped().Delete(&schemas.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}
