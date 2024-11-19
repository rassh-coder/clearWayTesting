package models

import "time"

type Asset struct {
	Name      string     `json:"name"`
	UID       uint       `json:"uId"`
	Data      *[]byte    `json:"-"`
	CreatedAt *time.Time `json:"createdAt"`
}
