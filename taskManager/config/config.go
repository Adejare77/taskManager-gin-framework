package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FullName  string `gorm:"column:fullName;not null"`
	Email     string `gorm:"column:email;unique;not null"`
	Password  string `gorm:"column:password;not null"`
	Tasks     []Task `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// Define Task Table
type Task struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"column:userID;not null;index"`
	TaskID    string `gorm:"column:taskID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Desc      string    `json:"description" gorm:"column:description;not null" binding:"required"`
	Title     string    `json:"title" gorm:"column:title;not null" binding:"required"`
	DueDate   time.Time `json:"dueDate" gorm:"column:dueDate;not null" binding:"required"`
	Status    string    `json:"status" gorm:"column:status;not null" binding:"required,status"`
	User      User      `gorm:"constraint:OnDelete:CASCADE;"`
}

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
	db.AutoMigrate(&User{}, &Task{})
	// db.AutoMigrate(&User{})

	DB = db
}
