package models

import "github.com/google/uuid"

type User struct {
	//User ID
	UserId uuid.UUID `json:"user_id"`
	//Login for User
	Login string `json:"login"`
	//Password for User
	Password string `json:"password"`
	//Roles that belong to the User
	Roles []string `json:"roles"`
}
