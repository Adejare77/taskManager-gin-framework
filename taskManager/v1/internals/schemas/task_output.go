package schemas

import "github.com/Adejare77/go/taskManager/internals/utilities"

type TaskOutput struct {
	TaskID    string `gorm:"column:taskID"`
	Desc      string `json:"description"`
	Title     string `json:"title"`
	StartDate string `json:"startDate" binding:"startdate"` // Default to time.Now()
	DueDate   string `json:"dueDate" binding:"duedate"`
	Status    string `json:"status"`
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

type PostTask struct {
	Desc      string `json:"description" binding:"required"`
	Title     string `json:"title" binding:"required"`
	StartDate string `json:"startDate" binding:"startdate"` // Default to time.Now()
	DueDate   string `json:"dueDate" binding:"required,duedate"`
	Status    string `json:"status"`
}

type DateTimeUpdate struct {
	StartDate string `json:"startDate" binding:"startdate"`
	DueDate   string `json:"dueDate"`
}

func ToTask(task PostTask) (Task, error) {
	startDate, err := utilities.StartTimeManipulator(task.StartDate)
	if err != nil {
		return Task{}, err
	}

	dueDate, err := utilities.DueTimeManipulator(task.DueDate, startDate)
	if err != nil {
		return Task{}, err
	}
	return Task{
		Desc:      task.Desc,
		Title:     task.Title,
		StartDate: startDate,
		DueDate:   dueDate,
		Status:    task.Status,
	}, nil
}
