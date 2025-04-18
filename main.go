package main

import (
	"log"

	"github.com/johnson-oragui/TaskFlow-Pro/api/database"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, continuing...")

	}

}

func main() {

	_, err := database.ConnectDatabase()
	if err != nil {
		log.Println(err.Error())
	}

}
