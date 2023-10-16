package model

import "github.com/google/uuid"

type UserLikeComment struct {
	User      User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    uuid.UUID `gorm:"null" json:"user_id"`
	Comment   Comment   `gorm:"ForeignKey:CommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CommentID string    `gorm:"null" json:"comment_id"`
}
