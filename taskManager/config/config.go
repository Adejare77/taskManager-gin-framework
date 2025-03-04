package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

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
		return fmt.Errorf("(DB Configuration) %v", err)
	}

	// Create Data Source Name (DSN)
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Taipei",
		cfg.User, cfg.Password, cfg.Host, cfg.DBName)

	// Open Database Connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("***ERROR DETECTED")
		return fmt.Errorf("(Connection) %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("(Connection Instance) %v", err)
	}

	maxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	connMaxLifetime, _ := time.ParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME"))

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxIdleTime(time.Duration(maxIdleConns))
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	// Auto Migrate tables into the Database
	if err := db.AutoMigrate(&schemas.User{}, &schemas.Task{}); err != nil {
		return fmt.Errorf("(Auto-Migration) %v", err)
	}

	DB = db
	fmt.Println("----------- db connected -------------")
	fmt.Println(DB)
	fmt.Println("----------- db connected -------------")
	return nil
}

// loadDBConfig loads database Configuration from .env file
func loadDBConfig() (*DBConfig, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("(DB_PORT) %v", err)
	}

	return &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}, nil
}
