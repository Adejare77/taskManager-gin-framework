package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Adejare77/go/taskManager/config"
	"github.com/Adejare77/go/taskManager/internals/handlers"
	"github.com/Adejare77/go/taskManager/internals/jobs"
	"github.com/Adejare77/go/taskManager/internals/middlewares"
	"github.com/Adejare77/go/taskManager/internals/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load Port if available
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000" // Defaults if not given
	}

	// Initialize application configuration
	if err := config.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Create Gin Instance
	app := gin.Default()

	// Use session middleware
	app.Use(sessions.Sessions("taskManager", config.SessionStore))

	// Start Cron Job
	if err := jobs.ScheduledStatusUpdater(); err != nil {
		handlers.Warning(fmt.Sprintf("error starting cron job: %s", err))
	}

	// Public Routes
	publicRoutes := app.Group("/")
	routes.PublicRoutes(publicRoutes)

	// Protected Routes
	protectedRoutes := app.Group("/", middlewares.AuthMiddleware())
	routes.ProtectedRoutes(protectedRoutes)

	// Health Check Endpoint
	app.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: app,
	}

	// Run server in gorouting
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not start the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	fmt.Println("Server shutting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
