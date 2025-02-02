package schemas

import (
	"time"

	"github.com/Adejare77/go/taskManager/internals/utilities"
	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FullName  string `gorm:"column:fullName;not null" binding:"required"`
	Email     string `gorm:"column:email;unique;not null" binding:"required,email"`
	Password  string `gorm:"column:password;not null" binding:"required"`
	Tasks     []Task `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// Hook to be called before saving a user
func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	user.Password, err = utilities.HashPassword(user.Password)
	return
}
