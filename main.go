package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis"
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
	RedisAddress := os.Getenv("REDIS_ADDRESS")
	RedisPassword := os.Getenv("REDIS_PASSWORD")
	RedisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		panic(err)
	}

	settings := mysql.ConnectionURL{
		Host:     MySQLHost,
		Database: MySQLDatabase,
		User:     MySQLUser,
		Password: MySQLPassword,
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisAddress,
		Password: RedisPassword,
		DB:       RedisDB,
	})

	db, err := db.Open(mysql.Adapter, settings)
	wb := NewWeb(db, rdb)

	log.Println("📳 Gummo server successfully started and listening on :8080")
	if err := http.ListenAndServe(":8080", wb); err != nil {
		panic(err)
	}
}
