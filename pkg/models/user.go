package models

import "time"

type User struct {
	ID           uint
	Login        string
	PasswordHash string
	CreatedAt    *time.Time
}

type UserInputs struct {
	Login    *string `json:"login,required"`
	Password *string `json:"password,required"`
}
