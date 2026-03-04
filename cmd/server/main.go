package main

import (
	"log"
	"net/http"

	"mini-analytics/internal/config"
	"mini-analytics/internal/repository"
)

func main() {
	cfg := config.Load()

	conn := repository.NewPostgresConn(cfg.DatabaseURL)
	defer conn.Close(nil)

	log.Println("Starting server on :8080")

	http.ListenAndServe(":8080", nil)
}
