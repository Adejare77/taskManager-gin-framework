package schemas

import (
	"time"

	"gorm.io/gorm"
)

// Define Task Table
type Task struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"column:userID;not null;index"`
	TaskID    string `gorm:"column:taskID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Desc      string    `json:"description" gorm:"column:description;not null"`
	Title     string    `json:"title" gorm:"column:title;not null"`
	StartDate time.Time `json:"startDate" gorm:"column:startDate;not null"`
	DueDate   time.Time `json:"dueDate" gorm:"column:dueDate;not null"`
	Status    string    `json:"status" gorm:"column:status;not null"`
	User      User      `gorm:"constraint:OnDelete:CASCADE;"`
}

// Hooks to be called before Creating a Task
func (task *Task) BeforeSave(tx *gorm.DB) error {
	if task.StartDate.Before(time.Now()) {
		task.Status = "in-progress"
	} else {
		task.Status = "pending"
	}
	return nil
}
