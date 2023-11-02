package model

import (
	"github.com/google/uuid"
	"time"
)

type UserVerify struct {
	UserID    uuid.UUID `gorm:"null" json:"user_id"`
	User      User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Verify    bool      `gorm:"default:false" json:"verify"`
	CreatedAt time.Time
}
