package models

import (
	"github.com/adejare77/taskmanager-gin-framework/config"
	"github.com/adejare77/taskmanager-gin-framework/internals/schemas"
)

func Create(user schemas.User) error {
	return config.DB.Create(&user).Error
}

func FindUserInfo(email string) (*schemas.User, error) {
	var user schemas.User
	if err := config.DB.Where("email=?", email).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(userID string) error {
	return config.DB.Delete(&schemas.User{}, "id = ?", userID).Error
}
