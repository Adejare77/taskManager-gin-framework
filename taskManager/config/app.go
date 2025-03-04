package config

import (
	"fmt"
	"log"

	"github.com/Adejare77/go/taskManager/internals/utilities"
	"github.com/joho/godotenv"
)

// Initialize sets up the application configuration
func Initialize() error {
	// load .env file
	if err := godotenv.Overload(); err != nil {
		log.Printf("Warning: No .env file found. Using environment variables")
	}

	// Initialize Session on Redis Server
	if err := InitSession(); err != nil {
		return fmt.Errorf("%v", err)
	}

	// Start Database Connection
	if err := Connect(); err != nil {
		return fmt.Errorf("%v", err)
	}

	// Register New Validators
	utilities.RegisterValidation()
	return nil
}
