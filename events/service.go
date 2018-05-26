package events

// Service contains business logic API related to events.
type Service interface {
	// Create  creates new event instance and returns ID.
	Create(Event) (int64, error)

	// One returns event instance by ID if it exists.
	One(int64) (Event, error)

	// All returns list of all events.
	All() []Event

	// Update updates existing events instance.
	Update(Event) error
}

var _ Service = (*eventsService)(nil)

type eventsService struct {
	repo Repository
}

// NewService instantiates new events service.
func NewService(repo Repository) Service {
	return eventsService{repo}
}

func (es eventsService) Create(event Event) (int64, error) {
	return es.repo.Create(event)
}

func (es eventsService) One(id int64) (Event, error) {
	return es.repo.One(id)
}

func (es eventsService) All() []Event {
	return es.repo.All()
}

func (es eventsService) Update(event Event) error {
	return es.repo.Update(event)
}
