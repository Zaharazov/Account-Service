package models

import "github.com/google/uuid"

type Student struct {
	//User ID
	UserId uuid.UUID `json:"user_id"`
	//Link to the photo for Student
	Photo string `json:"photo"`
	//Full Student name (last name, first name, patronymic)
	FullName string `json:"full_name"`
	//Student's group number
	Group string `json:"group"`
	//Record book number
	RecordBook string `json:"record_book"`
	//Description for Student
	Description string `json:"description"`
	//Student's e-mail
	Mail string `json:"mail"`
	//Link to Student's GitHub
	GitHub string `json:"github"`
}
