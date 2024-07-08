package models

import "github.com/google/uuid"

type Employer struct {
	UserId           uuid.UUID   `json:"user_id"`
	Photo            string      `json:"photo"`
	Name             string      `json:"name"`
	Focus            string      `json:"focus"`
	Description      string      `json:"description"`
	CreatedVacancies []uuid.UUID `json:"vacancies"`
}
