package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateEventHandler(service *EventService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var event Event

		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = service.CreateEvent(&event)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(event)
	}
}

func GetEventsHandler(service *EventService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
			return
		}

		events, err := service.GetEvents()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}

func EventByIDHandler(service *EventService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		idStr := r.URL.Query().Get("id")

		if idStr == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		switch r.Method {

		// GET EVENT BY ID
		case http.MethodGet:

			event, err := service.GetEventByID(id)

			if err == sql.ErrNoRows {
				http.Error(w, "event not found", http.StatusNotFound)
				return
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(event)

		// DELETE EVENT
		case http.MethodDelete:

			err := service.DeleteEvent(id)

			if err == sql.ErrNoRows {
				http.Error(w, "event not found", http.StatusNotFound)
				return
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status": "deleted",
			})

		// UPDATE EVENT
		case http.MethodPut:

			var event Event

			err := json.NewDecoder(r.Body).Decode(&event)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = service.UpdateEvent(id, &event)

			if err == sql.ErrNoRows {
				http.Error(w, "event not found", http.StatusNotFound)
				return
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(event)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
