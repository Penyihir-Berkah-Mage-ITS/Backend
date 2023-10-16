package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID
	Username       string `gorm:"unique" json:"username"`
	Phone          string `gorm:"unique" json:"phone"`
	Email          string `gorm:"unique" json:"email"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture"`
	AccountType    string `json:"account_type"`
	Gender         string `json:"gender"`
	IsVerified     bool   `gorm:"default:false" json:"is_verified"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterInput struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Gender         string `json:"gender"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profile_picture"`
}
