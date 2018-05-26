package postgres

import (
	"database/sql"
	"fmt"

	"github.com/anovakovic01/tap-backend/events"
)

var _ events.Repository = (*eventsRepository)(nil)

type eventsRepository struct {
	db *sql.DB
}

// NewEventsRepository instantiates a PostgreSQL implementation of events repository.
func NewEventsRepository(db *sql.DB) events.Repository {
	return eventsRepository{db}
}

func (er eventsRepository) Create(event events.Event) (int64, error) {
	q := `INSERT INTO events (owner, title, description, lat, lon, start, ending) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	var id int64
	err := er.db.
		QueryRow(q, event.Owner, event.Title, event.Description, event.Lat, event.Lon, event.Start, event.End).
		Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}

func (er eventsRepository) One(id int64) (events.Event, error) {
	q := `SELECT owner, title, description, lat, lon, start, ending FROM events WHERE id = $1`
	event := events.Event{ID: id}

	err := er.db.QueryRow(q, id).
		Scan(&event.Owner, &event.Title, &event.Description, &event.Lat, &event.Lon, &event.Start, &event.End)
	if err != nil {
		return events.Event{}, err
	}
	return event, nil
}

func (er eventsRepository) All() []events.Event {
	q := `SELECT id, owner, title, description, lat, lon, start, ending FROM events`

	rows, err := er.db.Query(q)
	if err != nil {
		fmt.Println(err)
		return []events.Event{}
	}
	defer rows.Close()

	items := []events.Event{}
	for rows.Next() {
		event := events.Event{}
		err := rows.
			Scan(&event.ID, &event.Owner, &event.Title, &event.Description, &event.Lat, &event.Lon, &event.Start, &event.End)
		if err != nil {
			fmt.Println(err)
			return []events.Event{}
		}
		items = append(items, event)
	}

	return items
}

func (er eventsRepository) Update(event events.Event) error {
	q := `UPDATE events
	      SET title = $1, description = $2, lat = $3, lon = $4, start = $5, ending = $6
		  WHERE id = $7 AND owner = $8`

	res, err := er.db.
		Exec(q, event.Title, event.Description, event.Lat, event.Lon, event.Start, event.End, event.ID, event.Owner)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return events.ErrNotFound
	}

	return nil
}
