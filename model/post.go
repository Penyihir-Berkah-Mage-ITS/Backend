package model

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	User       User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID     uuid.UUID `gorm:"null" json:"user_id"`
	Content    string    `json:"content"`
	Attachment string    `json:"attachment"`
	Likes      int64     `json:"likes"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Distance   string    `json:"distance"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostInput struct {
	Content    string `json:"content"`
	Attachment string `json:"attachment"`
}

type PostResponse struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	User         User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID       uuid.UUID `gorm:"null" json:"user_id"`
	Content      string    `json:"content"`
	Attachment   string    `json:"attachment"`
	Likes        int64     `json:"likes"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Distance     string    `json:"distance"`
	IsLiked      bool      `json:"is_liked"`
	TotalComment int64     `json:"total_comment"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
