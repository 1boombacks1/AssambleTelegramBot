package models

import "time"

type User struct {
	ChatID   int64     `json:"chat_id"`
	Username string    `json:"username"`
	Attempts byte      `json:"attempts"`
	Timer    time.Time `json:"timer"`
}

type MessageData struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int64 `json:"message_id"`
}

type AssambleInfo struct {
	ID                  int                      `json:"info_id"`
	AllUsersMessageData []MessageData            `json:"list_message_data"`
	ComeUsers           map[string][]MessageData `json:"come_users"`
	NotComeUsers        map[string][]MessageData `json:"not_come_users"`
}
