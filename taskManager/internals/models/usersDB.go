package models

import (
	"errors"
	"time"

	"github.com/Adejare77/go/taskManager/config"
	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FullName  string `json:"fullname" gorm:"column:fullName" binding:"required"`
	Email     string `json:"email" gorm:"column:email" binding:"required"`
	Password  string `json:"password" gorm:"column:password" binding:"required"`
	Tasks     []Task `json:"-" gorm:"column:tasks"`
}

var db = config.DB

func GetTasksByUserID(userID uint) ([]Task, error) {
	var user User
	if err := db.Preload("Tasks").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Tasks, nil
}

func GetTasksByStatus(userID uint, status string) ([]Task, error) {
	var user User
	if err := db.Preload("Tasks", "status = ?", status).First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Tasks, nil
}

func GetInfo(email string) (uint, string, error) {
	var user User
	err := db.Where("email=?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, "", err
	}
	return user.ID, user.Password, nil
}

func (user *User) Create() error {
	err := db.Create(&user).Error

	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(userID string) error {
	// db.Unscoped().Delete()  // for hard delete
	if err := db.Unscoped().Delete(&User{}, userID).Error; err != nil {
		return err
	}
	return nil
}
