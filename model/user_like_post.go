package model

import "github.com/google/uuid"

type UserLikePost struct {
	User   User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uuid.UUID `gorm:"null" json:"user_id"`
	Post   Post      `gorm:"ForeignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostID string    `gorm:"null" json:"post_id"`
}
