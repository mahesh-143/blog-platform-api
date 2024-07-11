package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Connection *sql.DB

func InitDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	host := os.Getenv("HOST")
	port := 5432
	user := "mahesh"
	password := "password"
	dbname := "blog"

	connectionStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	Connection, err = sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal("error validating database connection : ", err)
	}

	err = Connection.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database initialized")
}
