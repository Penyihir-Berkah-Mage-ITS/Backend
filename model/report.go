package model

import "github.com/google/uuid"

type Report struct {
	ID           uuid.UUID `gorm:"primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"null" json:"user_id"`
	User         User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Name         string    `gorm:"notNull" json:"name"`
	Address      string    `gorm:"notNull" json:"address"`
	Province     string    `gorm:"notNull" json:"province"`
	City         string    `gorm:"notNull" json:"city"`
	DetailReport string    `gorm:"notNull" json:"detail_report"`
	Proof        string    `json:"proof"`
}

type ReportInput struct {
	Name         string `gorm:"notNull" json:"name"`
	Address      string `gorm:"notNull" json:"address"`
	Province     string `gorm:"notNull" json:"province"`
	City         string `gorm:"notNull" json:"city"`
	DetailReport string `gorm:"notNull" json:"detail_report"`
	Proof        string `json:"proof"`
}
