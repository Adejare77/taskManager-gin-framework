package jobs

import (
	"github.com/Adejare77/go/taskManager/internals/models"
	"github.com/robfig/cron/v3"
)

func StatusUpdater() {
	// Create a cron instance
	c := cron.New()

	// schedule a task
	c.AddFunc("@every 60s", models.CheckStatus)

	// Start scheduler
	c.Start()
}
