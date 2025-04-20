package schemas

import (
	"time"

	"github.com/Adejare77/taskmanager/internals/utilities"
)

// Define Task Table
type Task struct {
	TaskID    string              `gorm:"type:uuid;default:gen_random_uuid();index"`
	UserID    string              `gorm:"not null;index;type:uuid"`
	Desc      string              `json:"description" gorm:"type:text;not null" binding:"required"`
	Title     string              `gorm:"type:text;not null" binding:"required"`
	StartDate *utilities.JSONTime `json:"start_date" gorm:"not null;type:timestamp" binding:"omitempty"`
	DueDate   utilities.JSONTime  `json:"due_date" gorm:"not null;type:timestamp" binding:"required"`
	Status    string              `gorm:"not null;default:pending"`
	User      User                `gorm:"foreignKey:UserID;" binding:"-"`
	CreatedAt time.Time           `gorm:"index"`
	UpdatedAt time.Time
}

type TaskQueryParams struct {
	Title  string `form:"title" binding:"omitempty"`
	Status string `form:"status" binding:"omitempty,oneof=pending in-progress completed"`
	Page   int    `form:"page" binding:"numeric,min=1"`
	Limit  int    `form:"limit" binding:"numeric,min=1"`
}

type TaskUriParam struct {
	TaskID string `uri:"task_id" binding:"required,uuid"`
}
