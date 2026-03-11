package db

import (
	"log"

	"golang-practice/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}
	var dsn = "host=localhost user=postgres password='Pass@123' dbname=userdb sslmode=disable"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("Connected to the database successfully.")

	DB.AutoMigrate(&models.User{})
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Failed to get database instance during close:", err)
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Println("Failed to close database:", err)
	}
}
