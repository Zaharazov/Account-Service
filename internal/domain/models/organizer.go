package models

import "github.com/google/uuid"

type Organizer struct {
	//User ID
	UserId uuid.UUID `json:"user_id"`
	//Link to the photo for Organizer
	Photo string `json:"photo"`
	//Organizer's name
	Name string `json:"name"`
	//Description for Organizer
	Description string `json:"description"`
	//Organizer-created events
	CreatedEvents []uuid.UUID `json:"events"`
}
