package main

import (
	"fmt"
	"log"

	"github.com/Adejare77/go/taskManager/config"
	"github.com/Adejare77/go/taskManager/internals/jobs"
	"github.com/Adejare77/go/taskManager/internals/middlewares"
	"github.com/Adejare77/go/taskManager/internals/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	app.Use(sessions.Sessions("taskManager", config.SessionStore))

	// Start Cron Job
	jobs.StatusUpdater()

	// Public Routes
	publicRoutes := app.Group("/")
	routes.PublicRoutes(publicRoutes)

	// Protected Routes
	protectedRoutes := app.Group("/", middlewares.AuthMiddleware())
	routes.ProtectedRoutes(protectedRoutes)

	fmt.Println("Running Server on Port 3000")
	if err := app.Run(":3000"); err != nil {
		log.Fatal("Could not Start the Server")
	}
}

// Handle 404 Not Found
// router.NoRoute(func(ctx *gin.Context) {
// 	ctx.JSON(http.StatusNotFound, gin.H{
// 		"error":   "Page not found",
// 		"message": "The requested URL was not found on this server.",
// 	})
// })
