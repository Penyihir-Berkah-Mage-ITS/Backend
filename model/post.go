package model

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID         string    `gorm:"unique" json:"id"`
	User       User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID     uuid.UUID `gorm:"null" json:"user_id"`
	Content    string    `json:"content"`
	Attachment string    `json:"attachment"`
	Like       int64     `json:"like"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
