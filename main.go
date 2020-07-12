package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"upper.io/db.v3"
	"upper.io/db.v3/mysql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, reading directly from env variables")
	}

	MySQLHost := os.Getenv("MYSQL_HOST")
	MySQLDatabase := os.Getenv("MYSQL_DATABASE")
	MySQLUser := os.Getenv("MYSQL_USER")
	MySQLPassword := os.Getenv("MYSQL_PASSWORD")

	settings := mysql.ConnectionURL{
		Host:     MySQLHost,
		Database: MySQLDatabase,
		User:     MySQLUser,
		Password: MySQLPassword,
	}

	sess, err := db.Open(mysql.Adapter, settings)
}
