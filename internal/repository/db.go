package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func NewPostgresConn(databaseURL string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Connected to PostgreSQL")
	return conn
}
