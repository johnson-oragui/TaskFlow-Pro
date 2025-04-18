package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	RedisUrl string
	DBURL    string
	Port     string
	AppUrl   string
}

func Load() *Config {
	HOST := os.Getenv("DB_HOST")
	USER := os.Getenv("DB_USERNAME")
	PASSWORD := os.Getenv("DB_PASSWORD")
	NAME := os.Getenv("DB_NAME")
	PORT := os.Getenv("DB_PORT")
	SSLMODE := os.Getenv("DB_SSLMODE")
	RedisUrl := os.Getenv("REDIS_URL")
	APPPORT := os.Getenv("PORT")
	APPURL := os.Getenv("APP_BASE_URL")

	if APPPORT == "" || APPURL == "" || RedisUrl == "" || HOST == "" || USER == "" || PASSWORD == "" || NAME == "" || PORT == "" || SSLMODE == "" {
		log.Fatalln("PORT REDIS_URL HOST , USER, PASSWORD, NAME, PORT, SSLMODE missing in configuration")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", HOST, USER, PASSWORD, NAME, PORT, SSLMODE)

	return &Config{
		RedisUrl: RedisUrl,
		DBURL:    dsn,
		Port:     APPPORT,
		AppUrl:   APPURL,
	}
}
