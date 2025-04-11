package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB


func ConnectDatabase() (*gorm.DB, error) {
	HOST := os.Getenv("DB_HOST")
	USER := os.Getenv("DB_USERNAME")
	PASSWORD := os.Getenv("DB_PASSWORD")
	NAME := os.Getenv("DB_NAME")
	PORT := os.Getenv("DB_PORT")
	SSLMODE := os.Getenv("DB_SSLMODE")
	if HOST == "" || USER == "" || PASSWORD == "" || NAME == "" || PORT == "" || SSLMODE == "" {
		log.Println("DATABASE HOST , USER, PASSWORD, NAME, PORT, SSLMODE missing in configuration.")
		return nil, fmt.Errorf("DATABASE HOST , USER, PASSWORD, NAME, PORT, SSLMODE missing in configuration")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Africa/West", HOST, USER, PASSWORD, NAME, PORT, SSLMODE)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil{
		return nil, err
	}

	DB = db
	log.Println("Connected to Database")
	return DB, nil

}
