package model

import "time"

type Message struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Content   string    `json:"content"`
	RoomID    string    `json:"roomId"`
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

func NewMessage(content, roomID, userID, username string) *Message {
	return &Message{
		Content:   content,
		RoomID:    roomID,
		UserID:    userID,
		Username:  username,
		Timestamp: time.Now(),
	}
}
