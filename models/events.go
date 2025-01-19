package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Event models
type Event struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	EventOwner     uuid.UUID `json:"event_owner"`
	EventStartTime time.Time `json:"event_start_time"`
	EventEndTime   time.Time `json:"event_end_time"`
	Participants   []string  `json:"participants"`
}

type EventRequest struct {
	Title         string   `json:"title" example:"Brainstorming meeting"`
	EventOwner    string   `json:"event_owner" example:"uuid"`
	EventTimeSlot string   `json:"event_time_slot" example:"02 Jan 2025 2-4 PM EST"`
	Participants  []string `json:"participants" example:"kevin,marco"`
}
