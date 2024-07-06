package models

type User struct {
	UserId   int32    `json:"user_id"`
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}
