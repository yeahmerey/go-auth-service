package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to DB")
}
