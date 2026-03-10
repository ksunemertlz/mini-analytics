package main

import (
	"fmt"
	"net/http"
)

func main() {

	db := InitDB()

	fmt.Println("Server started")

	http.HandleFunc("/event", CreateEventHandler(db))
	http.HandleFunc("/events", GetEventsHandler(db))
	http.HandleFunc("/event_by_id", EventByIDHandler(db))

	http.ListenAndServe(":9090", nil)
}
