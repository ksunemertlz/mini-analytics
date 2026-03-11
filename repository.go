package main

import "database/sql"

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event *Event) error {

	return r.db.QueryRow(
		"INSERT INTO event (user_id,event_type,page) VALUES ($1,$2,$3) RETURNING id,timestamp",
		event.UserID,
		event.EventType,
		event.Page,
	).Scan(&event.Id, &event.Timestamp)
}

func (r *EventRepository) GetAll() ([]Event, error) {

	rows, err := r.db.Query(
		"SELECT id,user_id,event_type,page,timestamp FROM event",
	)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

func (r *EventRepository) GetByID(id int) (*Event, error) {

	var event Event

	err := r.db.QueryRow(
		"SELECT id,user_id,event_type,page,timestamp FROM event WHERE id=$1",
		id,
	).Scan(
		&event.Id,
		&event.UserID,
		&event.EventType,
		&event.Page,
		&event.Timestamp,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventRepository) Delete(id int) error {

	res, err := r.db.Exec(
		"DELETE FROM event WHERE id=$1",
		id,
	)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *EventRepository) Update(id int, event *Event) error {

	return r.db.QueryRow(
		`UPDATE event
		 SET user_id=$2, event_type=$3, page=$4
		 WHERE id=$1
		 RETURNING id,user_id,event_type,page,timestamp`,
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
}
