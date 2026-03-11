package main

type EventService struct {
	repo *EventRepository
}

func NewEventService(repo *EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(event *Event) error {
	return s.repo.Create(event)
}

func (s *EventService) GetEvents() ([]Event, error) {
	return s.repo.GetAll()
}

func (s *EventService) GetEventByID(id int) (*Event, error) {
	return s.repo.GetByID(id)
}

func (s *EventService) DeleteEvent(id int) error {
	return s.repo.Delete(id)
}

func (s *EventService) UpdateEvent(id int, event *Event) error {
	return s.repo.Update(id, event)
}
