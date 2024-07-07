package models

import "github.com/google/uuid"

type User struct {
	UserId   uuid.UUID `json:"user_id"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	Roles    []string  `json:"roles"`
}
