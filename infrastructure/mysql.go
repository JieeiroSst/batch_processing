package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", username, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
	}
	return &Database{
		DB: db,
	}
}
