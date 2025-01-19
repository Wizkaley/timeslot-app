package models

import "github.com/gofrs/uuid"

type UserTimeSlotRequest struct {
	UserName  string   `json:"user_name" example:"eshan"`
	TimeSlots []string `json:"time_slots" example:"2 Jan 2025 2 - 4 PM EST,14 Jan 2025 6-9 PM EST"`
}

type User struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name" example:"eshan"`
}

type UserCreateRequest struct {
	Name string `json:"name" example:"eshan"`
}
