package main

type Event struct {
	Id        int    `json:"id"`
	UserID    int    `json:"user_id"`
	EventType string `json:"event_type"`
	Page      string `json:"page"`
	Timestamp string `json:"timestamp"`
}
