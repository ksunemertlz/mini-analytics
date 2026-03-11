package main

import "net/http"

func main() {

	db := InitDB()

	repo := NewEventRepository(db)

	service := NewEventService(repo)

	// CREATE
	http.Handle(
		"/event",
		LoggingMiddleware(CreateEventHandler(service)),
	)

	// GET ALL
	http.Handle(
		"/events",
		LoggingMiddleware(GetEventsHandler(service)),
	)

	// GET BY ID / DELETE / PUT
	http.Handle(
		"/event_by_id",
		LoggingMiddleware(EventByIDHandler(service)),
	)

	http.ListenAndServe(":9090", nil)
}
