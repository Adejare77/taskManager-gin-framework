package schemas

import "time"

type Task struct {
	TaskID    string `gorm:"column:taskID"`
	CreatedAt time.Time
	Desc      string `json:"description" gorm:"column:description"`
	Title     string `json:"title"`
	DueDate   string `json:"dueDate" gorm:"column:dueDate"`
	Status    string `json:"status"`
}
