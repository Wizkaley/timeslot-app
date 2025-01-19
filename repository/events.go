package repository

import (
	"timeslot-app/models"

	"github.com/jackc/pgx"
)

type EventRepoImplementation struct {
	db *pgx.Conn
}

func NewEventRepository(dbconn *pgx.Conn) EventRepo {
	return &EventRepoImplementation{
		db: dbconn,
	}
}

type EventRepo interface {
	CreateEvent(event models.Event) error
	GetEvent(eventID string) (models.Event, error)
	DeleteEvent(eventID string) error
	GetEventsForUser(username string) ([]models.Event, error)
}

func (er *EventRepoImplementation) CreateEvent(event models.Event) error {

	insertQuery := `INSERT INTO events (id, title, event_owner, event_start_time, event_end_time, participants) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := er.db.Exec(insertQuery, event.ID, event.Title, event.EventOwner, event.EventStartTime, event.EventEndTime, event.Participants)
	if err != nil {
		return err
	}
	return nil
}

func (er *EventRepoImplementation) GetEvent(eventID string) (models.Event, error) {

	var event models.Event
	err := er.db.QueryRow("SELECT * FROM events WHERE id = $1", eventID).Scan(&event.ID, &event.Title, &event.EventOwner, &event.EventStartTime, &event.EventEndTime, &event.Participants)
	if err != nil {
		return models.Event{}, err
	}
	return event, nil
}

func (er *EventRepoImplementation) DeleteEvent(eventID string) error {

	deleteQuery := `DELETE FROM events WHERE id = $1`
	_, err := er.db.Exec(deleteQuery, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (er *EventRepoImplementation) GetEventsForUser(username string) ([]models.Event, error) {
	qry := `select e.id, e.event_owner, e.title, e.event_start_time, e.event_end_time, e.participants from events e
		join users u on e.event_owner=u.id 
		where u.name = $1`

	rows, err := er.db.Query(qry, username)
	if err != nil {
		return []models.Event{}, err
	}
	var events []models.Event
	for rows.Next() {

		var event models.Event
		err := rows.Scan(&event.ID, &event.EventOwner, &event.Title, &event.EventStartTime, &event.EventEndTime, &event.Participants)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	return events, nil
}
