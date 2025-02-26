package config

import (
	"ToDo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := "host=localhost user=postgres password=obscurity dbname=todo_list port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&models.Task{}, &models.User{})
	DB = db
}
