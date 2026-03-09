package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Event struct {
	Id        int    `json:"id"`
	UserID    int    `json:"user_id"`
	EventType string `json:"event_type"`
	Page      string `json:"page"`
	Timestamp string `json:"timestamp"`
}

var events []Event

func main() {
	fmt.Println("Сервер запустился")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Кто-то зашёл на сайт")
		fmt.Fprintln(w, "Привет! Всё работает.")
	})

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {

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

		if event.EventType == "" {
			http.Error(w, "Error 400", http.StatusBadRequest)
			return
		}

		event.Id = len(events) + 1
		event.Timestamp = time.Now().Format(time.RFC3339)

		events = append(events, event)
		fmt.Println("Получили событие:", event)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(event)

	})

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	})

	http.HandleFunc("/event_by_id", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		var foundEvent *Event

		for i := range events {
			if fmt.Sprint(events[i].Id) == idStr {
				foundEvent = &events[i]
				break
			}
		}

		if foundEvent == nil {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(foundEvent)

	})

	http.HandleFunc("/event_by_id", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodDelete {
			http.Error(w, "only DELETE allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		deleted := false
		for i := range events {
			if fmt.Sprint(events[i].Id) == idStr {
				events = append(events[:i], events[i+1:]...)
				deleted = true
				break
			}
		}

		if !deleted {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	})

	http.ListenAndServe(":9090", nil)
}
