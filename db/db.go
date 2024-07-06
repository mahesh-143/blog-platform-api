package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var Connection *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "mahesh"
	password = "password"
	dbname   = "blog"
)

func InitDB() {

	connectionStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
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
