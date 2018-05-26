package events

import (
	"errors"
	"time"
)

// ErrNotFound indicates that event doesn't exist.
var ErrNotFound = errors.New("event not found")

// Repository contains API for reading, storing and editing
// events data.
type Repository interface {
	// Create creates new event instance and returns ID.
	Create(Event) (int64, error)

	// One returns event instance by ID if it exists.
	One(int64) (Event, error)

	// All returns list of all events stored in data store.
	All() []Event

	// Update updates existing events instance.
	Update(Event) error
}

// Event contains event related data.
type Event struct {
	ID          int64     `json:"id"`
	Owner       string    `json:"owner"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}
