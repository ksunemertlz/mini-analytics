package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {

	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=analytics sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	return db
}
