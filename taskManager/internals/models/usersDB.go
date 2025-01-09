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
}

var db = config.DB

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
	// if err := db.Delete(&User{}, "\"taskID\" = ?", userID).Error; err != nil {
	// 	return err
	// }

	return nil
}
