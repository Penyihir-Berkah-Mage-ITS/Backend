package model

import "time"

type Comment struct {
	ID         string `gorm:"primaryKey" json:"id"`
	User       User   `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID     string `gorm:"null" json:"user_id"`
	Post       Post   `gorm:"ForeignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostID     string `gorm:"null" json:"post_id"`
	Content    string `json:"content"`
	Attachment string `json:"attachment"`
	Like       int64  `json:"like"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
