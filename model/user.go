package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"unique" json:"username"`
	Phone          string    `gorm:"unique" json:"phone"`
	Email          string    `gorm:"unique" json:"email"`
	Password       string    `json:"password"`
	ProfilePicture string    `json:"profile_picture"`
	AccountType    string    `json:"account_type"`
	Gender         string    `json:"gender"`
	IsVerified     bool      `gorm:"default:false" json:"is_verified"`

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

type UserEditInput struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
}

type UserEditProfilePicture struct {
	ProfilePicture string `json:"profile_picture"`
}

type UserEditPassword struct {
	Password string `json:"password"`
}
