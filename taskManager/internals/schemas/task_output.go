package schemas

type TaskOutput struct {
	TaskID    string `gorm:"column:taskID"`
	Desc      string `json:"description" binding:"required"`
	Title     string `json:"title" binding:"required"`
	StartDate string `json:"startDate" binding:"date"` // Default to time.Now()
	DueDate   string `json:"dueDate" binding:"required,date"`
	Status    string `json:"status"`
}

func (task *TaskOutput) TableName() string {
	return "tasks"
}

func ToTaskOutput(task Task) TaskOutput {
	return TaskOutput{
		TaskID:    task.TaskID,
		Desc:      task.Desc,
		Title:     task.Title,
		StartDate: task.StartDate.Format("2006-01-02 15:04:05"),
		DueDate:   task.DueDate.Format("2006-01-02 15:04:05"),
		Status:    task.Status,
	}
}
