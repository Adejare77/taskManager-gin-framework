package config

import (
	"fmt"

	"github.com/Adejare77/taskmanager/internals/handlers"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Overload(); err != nil {
		handlers.Warning("No .env file found. Using environment variables")
	}
}

func Initialize() error {
	if err := Connect(); err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := InitSession(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
