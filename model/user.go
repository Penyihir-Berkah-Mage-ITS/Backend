package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID
	Username   string `gorm:"unique" json:"username"`
	Phone      string
	Email      string
	Password   string
	Gender     string
	IsVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserLogin struct {
	Username string
	Password string
}

type UserRegister struct {
	Username string
	Password string
	Gender   string
	Email    string
}
