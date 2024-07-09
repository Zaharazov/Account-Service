package models

import "github.com/google/uuid"

type Employer struct {
	//User ID
	UserId uuid.UUID `json:"user_id"`
	//Link to the photo for Employer
	Photo string `json:"photo"`
	//Name of the organization
	Name string `json:"name"`
	//Scope of the organization's work
	Focus string `json:"focus"`
	//Description for Employer
	Description string `json:"description"`
	//Employer-created vacancies
	CreatedVacancies []uuid.UUID `json:"vacancies"`
}
