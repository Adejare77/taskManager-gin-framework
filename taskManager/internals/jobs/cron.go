package jobs

import (
	"fmt"
	"log"
	"os"

	"github.com/Adejare77/go/taskManager/internals/models"
	"github.com/robfig/cron/v3"
)

// func StatusUpdater() {
// 	// Create a cron instance
// 	c := cron.New()

// 	// schedule a task
// 	c.AddFunc("@every 60s", models.CheckStatus)

// 	// Start scheduler
// 	c.Start()
// }

var cronScheduler *cron.Cron

func StatusUpdater() error {
	// Load cron schedule from environment variable
	schedule := os.Getenv("CRON_SCHEDULE")
	if schedule == "" {
		schedule = "@60s" // Default schedule
	}

	// Create a new cron instance or scheduler
	cronScheduler = cron.New()

	// Add the task status update job
	_, err := cronScheduler.AddFunc(schedule, updateTaskStatus)
	if err != nil {
		return fmt.Errorf("failed to add cron job: %v", err)
	}

	// start the cron scheduler
	cronScheduler.Start()

	log.Println("Cron jobs started")
	return nil
}

// StopCronJobs stops the cron jobs gracefully
func updateTaskStatus() {
	log.Println("Updating task status...")
	models.UpdateTaskStatus()
	log.Println("Task status updated successfully")
}
