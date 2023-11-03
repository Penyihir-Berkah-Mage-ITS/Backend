package model

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID      string    `gorm:"primaryKey" json:"id"`
	User    User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID  uuid.UUID `gorm:"null" json:"user_id"`
	Post    Post      `gorm:"ForeignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostID  string    `gorm:"null" json:"post_id"`
	Content string    `json:"content"`
	Like    int64     `json:"like"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentResponse struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	User         User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID       uuid.UUID `gorm:"null" json:"user_id"`
	Post         Post      `gorm:"ForeignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostID       string    `gorm:"null" json:"post_id"`
	Content      string    `json:"content"`
	Like         int64     `json:"like"`
	IsLiked      bool      `json:"is_liked"`
	TotalComment int64     `json:"total_comment"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentInput struct {
	Content string `json:"content"`
}
