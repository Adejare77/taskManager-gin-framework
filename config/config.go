package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Adejare77/go/taskManager/internals/handlers"
	"github.com/Adejare77/go/taskManager/internals/schemas"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

var DB *gorm.DB

func Connect() error {
	// Load database Configuration
	cfg, err := loadDBConfig()
	if err != nil {
		return fmt.Errorf("DB Configuration Error %v", err)
	}

	// Create Data Source Name (DSN)
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=5432 sslmode=disable TimeZone=Africa/Lagos",
		cfg.User, cfg.Password, cfg.Host, cfg.DBName)

	// Open Database Connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("connection Instance %v", err)
	}

	maxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	connMaxLifetime, err := time.ParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME"))
	if err != nil {
		connMaxLifetime = time.Duration(30)
		handlers.Warning("invalid DB_CONN_MAX_LIFETIME. Defaults to 30 minutes")
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxIdleTime(time.Duration(maxIdleConns))
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	// Auto Migrate tables into the Database
	if err := db.AutoMigrate(&schemas.User{}, &schemas.Task{}); err != nil {
		return fmt.Errorf("Auto-Migration Error %v", err)
	}

	DB = db
	return nil
}

// loadDBConfig loads database Configuration from .env file
func loadDBConfig() (*DBConfig, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 5432
		handlers.Warning("invalid DB_PORT. Defaults to 5432")
		// return nil, fmt.Errorf("(DB_PORT) %v", err)
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
		handlers.Warning("invalid HOST. Defaults to `localhost`")
	}

	return &DBConfig{
		Host:     host,
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}, nil
}
