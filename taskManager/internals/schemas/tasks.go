package schemas

import (
	"errors"
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
	DueDate   time.Time `json:"dueDate" gorm:"column:dueDate;not null"`
	StartDate time.Time `json:"startDate" gorm:"column:startDate;not null"`
	Status    string    `json:"status" gorm:"column:status;not null"`
	User      User      `gorm:"constraint:OnDelete:CASCADE;"`
}

// Hooks to be called before Creating a Task
func (task *Task) BeforeSave(tx *gorm.DB) (err error) {
	if task.StartDate.IsZero() {
		task.StartDate = time.Now()
		task.Status = "in-progress"
		return
	} else if task.StartDate.Before(time.Now()) {
		err = errors.New("start date cannot be in the past")
		return
	} else if task.StartDate.After(task.DueDate) {
		err = errors.New("start date must be before dueDate")
		return
	}
	task.Status = "pending"
	return
}
