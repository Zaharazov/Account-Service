package models

import "github.com/google/uuid"

type Student struct {
	UserId      uuid.UUID `json:"user_id"`
	Photo       string    `json:"photo"`
	FullName    string    `json:"full_name"`
	Group       string    `json:"group"`
	RecordBook  string    `json:"record_book"`
	Description string    `json:"description"`
	Mail        string    `json:"mail"`
	GitHub      string    `json:"github"`
}
