package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type RecommendSlotsResponse struct {
	MatchedSlots []TimeSlotStartAndEnd `json:"Matched Slots"`
	PartialSlots []MatchingEventSlots  `json:"Partially Matched Slots"`
}

type MatchingEventSlots struct {
	Slot                    TimeSlotStartAndEnd
	AvailableParticipants   []string `json:"Available Participants"`
	UnavailableParticipants []string `json:"Unavailable Participants"`
}

type TimeSlotStartAndEnd struct {
	StartTime time.Time `json:"Start Time"`
	EndTime   time.Time `json:"End Time"`
}

type RecommendSlotsRequest struct {
	Organizer     string   `json:"organizer" example:"eshan"`
	Participants  []string `json:"participants" example:"kevin,marco"`
	EventDuration int      `json:"event_duration" example:"1"`
}

type Participant struct {
	Name      string
	TimeSlots []TimeSlotStartAndEnd
}

type TimeSlot struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	TimeSlot string    `json:"time_slot"`
}

type TimeSlotResponse struct {
	UserName  string   `json:"user_name"`
	TimeSlots []string `json:"time_slot"`
}

type ServiceMessage struct {
	Message string `json:"message"`
}
