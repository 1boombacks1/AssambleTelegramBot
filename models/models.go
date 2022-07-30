package models

import "time"

type User struct {
	ChatID   int64     `json:"chat_id"`
	Username string    `json:"username"`
	Attempts byte      `json:"attempts"`
	Timer    time.Time `json:"timer"`
}
