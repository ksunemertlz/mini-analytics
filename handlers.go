package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateEventHandler(db *sql.DB) http.HandlerFunc {

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

		err = db.QueryRow(
			"INSERT INTO event (user_id,event_type,page) VALUES ($1,$2,$3) RETURNING id,timestamp",
			event.UserID,
			event.EventType,
			event.Page,
		).Scan(&event.Id, &event.Timestamp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(event)

	}
}

func GetEventsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query("SELECT id,user_id,event_type,page,timestamp FROM event")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var events []Event

		for rows.Next() {

			var e Event

			err := rows.Scan(
				&e.Id,
				&e.UserID,
				&e.EventType,
				&e.Page,
				&e.Timestamp,
			)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			events = append(events, e)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)

	}
}

func EventByIDHandler(db *sql.DB) http.HandlerFunc {
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

		case http.MethodGet:

			var event Event

			err := db.QueryRow(
				"SELECT id,user_id,event_type,page,timestamp FROM events WHERE id=$1",
				id,
			).Scan(
				&event.Id,
				&event.UserID,
				&event.EventType,
				&event.Page,
				&event.Timestamp,
			)

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

		case http.MethodDelete:

			res, err := db.Exec(
				"DELETE FROM event WHERE id=$1",
				id,
			)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			count, _ := res.RowsAffected()

			if count == 0 {
				http.Error(w, "event not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status": "deleted",
			})

		case http.MethodPut:

			var event Event

			// читаем JSON из запроса
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = db.QueryRow(
				"UPDATE event SET user_id=$2, event_type=$3, page=$4 WHERE id=$1 RETURNING id, user_id, event_type, page, timestamp",
				id,
				event.UserID,
				event.EventType,
				event.Page,
			).Scan(
				&event.Id,
				&event.UserID,
				&event.EventType,
				&event.Page,
				&event.Timestamp,
			)

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
