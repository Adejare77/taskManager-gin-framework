package config

import (
	"fmt"
	"log"
	"time"

	"github.com/Adejare77/go/taskManager/internals/schemas"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(user string, password string, dbname string) {
	// Data Source Name (dns)
	dsn := fmt.Sprintf("user=%s password=%s host=localhost dbname=%s port=5432 sslmode=disable TimeZone=Asia/Taipei",
		user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to Connect to Database")
	}

	// Other settings;
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(5) // concurrent connection
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// Auto Migrate tables into the Database
	db.AutoMigrate(&schemas.User{}, &schemas.Task{})

	DB = db
}
