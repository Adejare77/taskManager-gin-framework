package schemas

import (
	"time"

	"github.com/Adejare77/taskmanager/internals/utilities"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FullName  string `gorm:"column:fullName;not null" binding:"required"`
	Email     string `gorm:"email;unique;not null" binding:"required,email"`
	Password  string `gorm:"not null" binding:"required"`
	Tasks     []Task `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Hook to be called before saving a user
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Password, err = utilities.HashPassword(user.Password)
	return
}
