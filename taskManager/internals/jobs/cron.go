package jobs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Adejare77/go/taskManager/internals/handlers"
	"github.com/Adejare77/go/taskManager/internals/models"
	"github.com/robfig/cron/v3"
)

func ScheduledStatusUpdater() error {
	// Load cron schedule from environment variable
	// schedule, err := time.ParseDuration(os.Getenv("CRON_SCHEDULE"))
	schedule, err := strconv.Atoi(os.Getenv("CRON_SCHEDULE"))
	if err != nil {
		handlers.Warning("Invalid cron schedule time. Default to 60s")
		schedule = 60
	}

	// Create a new cron instance or scheduler
	cronScheduler := cron.New()
	// Add the task status update job
	cronID, err := cronScheduler.AddFunc(fmt.Sprintf("@every %ds", schedule), models.StatusUpdater)
	if err != nil {
		return fmt.Errorf("failed to add cron job: %v", err)
	}

	// start the cron scheduler
	cronScheduler.Start()

	handlers.Info(map[string]any{
		"cron_job ID": cronID,
	}, "cron job started..")

	return nil
}
