package config

import (
	"os"

	"github.com/Adejare77/go/taskManager/internals/utilities"
	"github.com/joho/godotenv"
)

func init() {
	// load .env file
	godotenv.Overload()

	addr := os.Getenv("REDIS_ADDRESS")
	pwd := os.Getenv("PASSWORD")
	db_name := os.Getenv("DB")
	user := os.Getenv("USER")
	secretKey := os.Getenv("SECRET_KEY")

	// Start Database Connection
	Connect(user, pwd, db_name)

	// Create Session on Redis Server
	InitSession(addr, pwd, secretKey)

	// Register New Validators
	utilities.RegisterValidation()

}
