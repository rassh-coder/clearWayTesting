package models

import "time"

type Session struct {
	ID        string     `json:"id"`
	UID       uint       `json:"uId"`
	CreatedAt *time.Time `json:"createdAt"`
	Ip        string     `json:"ip"`
	IsActive  bool       `json:"isActive"`
}
