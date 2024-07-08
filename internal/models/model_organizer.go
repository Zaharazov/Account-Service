package models

import "github.com/google/uuid"

type Organizer struct {
	UserId        uuid.UUID   `json:"user_id"`
	Photo         string      `json:"photo"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	CreatedEvents []uuid.UUID `json:"events"`
}
